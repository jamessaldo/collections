package persistence

import (
	"authorization/domain"
	"authorization/service"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type Seed struct {
	db *gorm.DB
}

type EndpointYAML struct {
	Endpoints []struct {
		ID     ulid.ULID `yaml:"id"`
		Name   string    `yaml:"name"`
		Path   string    `yaml:"path"`
		Method string    `yaml:"method"`
	} `yaml:"endpoints"`
}

type RoleYAML struct {
	ID        ulid.ULID       `yaml:"id"`
	Name      domain.RoleType `yaml:"name"`
	Endpoints []struct {
		Name string `yaml:"name"`
	} `yaml:"endpoints"`
}

// Execute will executes the given seeder method
func Execute(db *gorm.DB, seedMethodNames ...string) {
	s := Seed{db}

	seedType := reflect.TypeOf(s)

	// Execute all seeders if no method name is given
	if len(seedMethodNames) == 0 {
		log.Info().Msg("Seeding database...")
		// We are looping over the method on a Seed struct
		for i := 0; i < seedType.NumMethod(); i++ {
			// Get the method in the current iteration
			method := seedType.Method(i)
			// Execute seeder
			seed(s, method.Name)
		}
	}

	// Execute only the given method names
	for _, item := range seedMethodNames {
		seed(s, item)
	}
}

func seed(s Seed, seedMethodName string) {
	start := time.Now()
	// Get the reflect value of the method
	m := reflect.ValueOf(s).MethodByName(seedMethodName)
	// Exit if the method doesn't exist
	if !m.IsValid() {
		log.Fatal().Err(fmt.Errorf("method %s does not exist", seedMethodName)).Msg("Failed to seed database")
	}
	// Execute the method
	log.Info().Str("method name", seedMethodName).Msg("Seeding database...")
	m.Call(nil)
	duration := time.Since(start)
	log.Info().Int("duration", int(math.Ceil(duration.Seconds()))).Msg("Successfully seeded database")
}

func generateUsers(jobs chan<- domain.User) {
	now := time.Now()
	for i := 0; i < 1_111_111; i++ {
		userId := uuid.NewV4()
		user_data := domain.User{
			ID:        userId,
			FirstName: faker.ChineseFirstName(),
			LastName:  faker.ChineseLastName(),
			Email:     faker.Email(),
			Username:  faker.Username(),
			Password:  "",
			AvatarURL: "",
			Provider:  "Google",
			Verified:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		now = now.AddDate(0, 0, 10)
		jobs <- user_data
	}
}

func dispatchWorkers(db *gorm.DB, jobs <-chan domain.User, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex < 100; workerIndex++ {
		wg.Add(1)
		go func(workerIndex int, db *gorm.DB, jobs <-chan domain.User, wg *sync.WaitGroup) {
			counter := 0
			var users []domain.User
			for {
				if counter < 15000 {
					job, more := <-jobs
					if !more {
						break
					}
					users = append(users, job)
					counter++
					continue
				}
				break
			}
			userSeedBatchRoutine(db, users, counter, workerIndex)
			wg.Done()
		}(workerIndex, db, jobs, wg)
		time.Sleep(25 * time.Millisecond)
	}
}

func (s Seed) UserSeed() {
	maxJobBuffer := 100

	jobs := make(chan domain.User, maxJobBuffer)
	wg := new(sync.WaitGroup)

	go dispatchWorkers(s.db, jobs, wg)
	generateUsers(jobs)

	close(jobs)
	wg.Wait()
}

func userSeedBatchRoutine(db *gorm.DB, user_list []domain.User, counter, workerIndex int) {
	uow, err := service.NewUnitOfWork(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create unit of work")
	}

	userErr := uow.User.AddBatch(user_list)
	if userErr != nil {
		log.Error().Err(err).Msg("Failed to insert users")
	}

	log.Info().Msg(fmt.Sprintf("=> worker %d inserted %d data", workerIndex, counter))
}

func _(db *gorm.DB, user_data *domain.User, counter, workerIndex int) {
	uow, err := service.NewUnitOfWork(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create unit of work")
	}

	tx, err := uow.Begin(&gorm.Session{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to begin transaction")
	}

	defer func() {
		tx.Rollback()
	}()

	_, userErr := uow.User.Add(user_data, tx)
	if userErr != nil {
		log.Error().Err(err).Msg("Failed to insert user")
		tx.Rollback()
	}

	tx.Commit()

	log.Info().Msg(fmt.Sprintf("=> worker %d inserted %d data", workerIndex, counter))
}

func (s Seed) AccessSeed() {
	uow, err := service.NewUnitOfWork(s.db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create unit of work")
	}

	tx, err := uow.Begin(&gorm.Session{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to begin transaction")
	}

	defer func() {
		tx.Rollback()
	}()

	endpointDatas := readYAML("endpoints.yml")
	var endpointYAML EndpointYAML
	err = yaml.Unmarshal(endpointDatas, &endpointYAML)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal endpoint data")
	}

	roleDatas := readYAML("roles.yml")
	var roleYAML []RoleYAML
	err = yaml.Unmarshal(roleDatas, &roleYAML)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal role data")
	}

	cacheEndpoint := make(map[string]*domain.Endpoint)

	for _, endpoint := range endpointYAML.Endpoints {
		endpointData := domain.NewEndpoint(endpoint.ID, endpoint.Name, endpoint.Path, endpoint.Method)
		cacheEndpoint[endpoint.Name] = endpointData
	}

	for _, role := range roleYAML {
		roleData := domain.NewRole(role.ID, role.Name)
		for _, endpoint := range role.Endpoints {
			if val, ok := cacheEndpoint[endpoint.Name]; ok {
				roleData.AddEndpoints(val)
			}
		}
		_, roleErr := uow.Role.Add(roleData, tx)
		if roleErr != nil {
			log.Error().Err(roleErr).Msg("Failed to insert role")
		}
	}

	tx.Commit()
}

func readYAML(filename string) []byte {
	filePath := filepath.Join("data", filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("Failed to read role data %s", filename))
	}
	return data
}
