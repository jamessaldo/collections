package persistence

import (
	"authorization/domain/model"
	"authorization/service"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
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
		ID     uuid.UUID `yaml:"id"`
		Name   string    `yaml:"name"`
		Path   string    `yaml:"path"`
		Method string    `yaml:"method"`
	} `yaml:"endpoints"`
}

type RoleYAML struct {
	ID        uuid.UUID      `yaml:"id"`
	Name      model.RoleType `yaml:"name"`
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

func generateUsers(jobs chan<- model.User) {
	now := time.Now()
	for i := 0; i < 1_111_111; i++ {
		userId := uuid.NewV4()
		user_data := model.User{
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

func dispatchWorkers(db *gorm.DB, jobs <-chan model.User, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex < 100; workerIndex++ {
		wg.Add(1)
		go func(workerIndex int, db *gorm.DB, jobs <-chan model.User, wg *sync.WaitGroup) {
			counter := 0
			var users []model.User
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

	jobs := make(chan model.User, maxJobBuffer)
	wg := new(sync.WaitGroup)

	go dispatchWorkers(s.db, jobs, wg)
	generateUsers(jobs)

	close(jobs)
	wg.Wait()
}

func userSeedBatchRoutine(db *gorm.DB, user_list []model.User, counter, workerIndex int) {
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

func _(db *gorm.DB, user_data *model.User, counter, workerIndex int) {
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

	endpointPath := filepath.Join("data", "endpoints.yml")
	endpointDatas, err := os.ReadFile(endpointPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ")
	}

	var endpointYAML EndpointYAML
	err = yaml.Unmarshal(endpointDatas, &endpointYAML)
	if err != nil {
		panic(err)
	}

	rolePath := filepath.Join("data", "roles.yml")
	roleDatas, err := os.ReadFile(rolePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ")
	}

	var roleYAML []RoleYAML
	err = yaml.Unmarshal(roleDatas, &roleYAML)
	if err != nil {
		panic(err)
	}

	cacheEndpoint := make(map[string]*model.Endpoint)

	for _, endpoint := range endpointYAML.Endpoints {
		endpointData := model.Endpoint{
			ID:     endpoint.ID,
			Name:   endpoint.Name,
			Path:   endpoint.Path,
			Method: endpoint.Method,
		}
		cacheEndpoint[endpoint.Name] = &endpointData
	}

	for _, role := range roleYAML {
		roleData := model.Role{
			ID:   role.ID,
			Name: role.Name,
		}
		for _, endpoint := range role.Endpoints {
			roleData.AddEndpoints(cacheEndpoint[endpoint.Name])
		}
		_, roleErr := uow.Role.Add(&roleData, tx)
		if roleErr != nil {
			log.Error().Err(roleErr).Msg("Failed to insert role")
		}
	}

	tx.Commit()
}
