package integration

import (
	"authorization/config"
	"authorization/domain/model"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/service/handlers"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

var (
	Db            *gorm.DB
	cleanupDocker func()
	Bus           *service.MessageBus
)

var _ = BeforeSuite(func() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	// setup *gorm.Db with docker
	Db, cleanupDocker = setupGormWithDocker()
})

var _ = AfterSuite(func() {
	// cleanup resource
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	err := Db.Exec(`DROP SCHEMA public CASCADE;`).Error
	Ω(err).To(Succeed())
	err = Db.Exec(`CREATE SCHEMA public;`).Error
	Ω(err).To(Succeed())

	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..")
	err = os.Chdir(dir)
	if err != nil {
		panic(err)
	}

	err = Db.AutoMigrate(
		&model.User{},
		&model.Endpoint{},
		&model.Role{},
		&model.Access{},
		&model.Membership{},
		&model.Team{},
		&model.Invitation{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	uow, err := service.NewUnitOfWork(Db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create unit of work")
	}

	asynqClient := worker.CreateAsynqClientMock()

	mailer := worker.NewMailerMock(asynqClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create mailer")
	}

	Ω(err).To(Succeed())

	Bus = service.NewMessageBus(handlers.COMMAND_HANDLERS, uow, mailer)

	persistence.Execute(Db, "AccessSeed")
})

func setupGormWithDocker() (*gorm.DB, func()) {
	pool, err := dockertest.NewPool("")
	chk(err)

	runDockerOpt := &dockertest.RunOptions{
		Repository: "postgres", // image
		Tag:        "14",       // version
		Env:        []string{"POSTGRES_PASSWORD=" + config.StorageConfig.DBPassword, "POSTGRES_DB=" + config.StorageConfig.DBName},
	}

	fnConfig := func(config *docker.HostConfig) {
		config.AutoRemove = true                     // set AutoRemove to true so that stopped container goes away by itself
		config.RestartPolicy = docker.NeverRestart() // don't restart container
	}

	resource, err := pool.RunWithOptions(runDockerOpt, fnConfig)
	chk(err)
	// call clean up function to release resource
	fnCleanup := func() {
		err := resource.Close()
		chk(err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	dsn := fmt.Sprintf("%s://%s:%s@%s/%s", config.StorageConfig.DBDriver, config.StorageConfig.DBUser, config.StorageConfig.DBPassword, hostAndPort, config.StorageConfig.DBName)

	var gdb *gorm.DB
	// retry until db server is ready
	err = pool.Retry(func() error {
		gdb, err = gorm.Open(postgres.New(postgres.Config{
			DriverName: "pgx",
			DSN:        dsn,
		}), &gorm.Config{
			PrepareStmt: true,
			Logger:      logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return err
		}

		db, err := gdb.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	})
	chk(err)

	// container is ready, return *gorm.Db for testing
	return gdb, fnCleanup
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
