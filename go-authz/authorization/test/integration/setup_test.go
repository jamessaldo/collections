package integration

import (
	"authorization/config"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/service/handlers"
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func TestDocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Docker Suite")
}

var (
	Pool          *pgxpool.Pool
	cleanupDocker func()
	Bus           *service.MessageBus
)

var _ = BeforeSuite(func() {
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	// setup *pgxpool.Pool with docker
	Pool, cleanupDocker = setupPoolWithDocker()
})

var _ = AfterSuite(func() {
	// cleanup resource
	cleanupDocker()
})

var _ = BeforeEach(func() {
	// clear db tables before each test
	ctx := context.Background()
	_, err := Pool.Exec(ctx, `DROP SCHEMA public CASCADE;`)
	Ω(err).To(Succeed())
	_, err = Pool.Exec(ctx, `CREATE SCHEMA public;`)
	Ω(err).To(Succeed())

	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..", "..")
	err = os.Chdir(dir)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to change directory")
	}

	persistence.Migration(Pool)

	uow, err := service.NewUnitOfWork(Pool)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create unit of work")
	}

	asynqClient := worker.CreateAsynqClientMock()

	mailer := worker.NewMailerMock(asynqClient)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create mailer")
	}

	Ω(err).To(Succeed())

	Bus = service.NewMessageBus(handlers.COMMAND_HANDLERS, uow, mailer)

	persistence.Execute(Pool, "AccessSeed")
})

func setupPoolWithDocker() (*pgxpool.Pool, func()) {
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

	var gpool *pgxpool.Pool
	// retry until db server is ready
	err = pool.Retry(func() error {
		ctx := context.Background()
		gpool, err = pgxpool.New(ctx, dsn)
		if err != nil {
			return err
		}

		return gpool.Ping(ctx)
	})
	chk(err)

	// container is ready, return *pgxpool.Pool for testing
	return gpool, fnCleanup
}

func chk(err error) {
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to setup docker")
	}
}
