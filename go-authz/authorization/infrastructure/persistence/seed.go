package persistence

import (
	"authorization/domain"
	"authorization/service"
	"authorization/util"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"gopkg.in/yaml.v3"
)

type Seed struct {
	pool *pgxpool.Pool
}

type EndpointYAML struct {
	Endpoints []struct {
		Name   string `yaml:"name"`
		Path   string `yaml:"path"`
		Method string `yaml:"method"`
	} `yaml:"endpoints"`
}

type RoleYAML struct {
	Name      domain.RoleType `yaml:"name"`
	Endpoints []struct {
		Name string `yaml:"name"`
	} `yaml:"endpoints"`
}

// Execute will executes the given seeder method
func Execute(pool *pgxpool.Pool, seedMethodNames ...string) {
	s := Seed{pool}

	seedType := reflect.TypeOf(s)

	// Execute all seeders if no method name is given
	if len(seedMethodNames) == 0 {
		log.Info().Caller().Msg("Seeding database...")
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
	start := util.GetTimestampUTC()
	// Get the reflect value of the method
	m := reflect.ValueOf(s).MethodByName(seedMethodName)
	// Exit if the method doesn't exist
	if !m.IsValid() {
		log.Fatal().Caller().Err(fmt.Errorf("method %s does not exist", seedMethodName)).Msg("Failed to seed database")
	}
	// Execute the method
	log.Info().Caller().Str("method name", seedMethodName).Msg("Seeding database...")
	m.Call(nil)
	duration := time.Since(start)
	log.Info().Caller().Int("duration", int(math.Ceil(duration.Seconds()))).Msg("Successfully seeded database")
}

func generateUsers(jobs chan<- domain.User) {
	for i := 0; i < 1_111_111; i++ {
		user_data := domain.NewUser(faker.ChineseFirstName(), faker.ChineseLastName(), faker.Email(), "", "Google", true)
		jobs <- user_data
	}
}

func dispatchWorkers(pool *pgxpool.Pool, jobs <-chan domain.User, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex < 100; workerIndex++ {
		wg.Add(1)
		go func(workerIndex int, pool *pgxpool.Pool, jobs <-chan domain.User, wg *sync.WaitGroup) {
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
			userSeedBatchRoutine(pool, users, counter, workerIndex)
			wg.Done()
		}(workerIndex, pool, jobs, wg)
		time.Sleep(25 * time.Millisecond)
	}
}

func (s Seed) UserSeed() {
	maxJobBuffer := 100

	jobs := make(chan domain.User, maxJobBuffer)
	wg := new(sync.WaitGroup)

	go dispatchWorkers(s.pool, jobs, wg)
	generateUsers(jobs)

	close(jobs)
	wg.Wait()
}

func userSeedBatchRoutine(pool *pgxpool.Pool, user_list domain.Users, counter, workerIndex int) {
	uow, err := service.NewUnitOfWork(pool)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create unit of work")
	}
	ctx := context.Background()
	tx, err := uow.Begin(ctx)

	// for _, user := range user_list {
	// 	_, userErr := uow.User.Add(user, tx)
	// 	if userErr != nil {
	// 		log.Error().Caller().Err(err).Msg("Failed to insert users")
	// 	}
	// }
	userErr := uow.User.AddBatch(ctx, user_list, tx)
	if userErr != nil {
		log.Error().Caller().Err(err).Msg("Failed to insert users")
	}
	tx.Commit(ctx)

	log.Info().Caller().Msg(fmt.Sprintf("=> worker %d inserted %d data", workerIndex, counter))
}

func _(pool *pgxpool.Pool, user_data domain.User, counter, workerIndex int) {
	uow, err := service.NewUnitOfWork(pool)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create unit of work")
	}

	ctx := context.Background()
	tx, err := uow.Begin(ctx)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to begin transaction")
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	_, userErr := uow.User.Add(ctx, user_data, tx)
	if userErr != nil {
		log.Error().Caller().Err(err).Msg("Failed to insert user")
		tx.Rollback(ctx)
	}

	tx.Commit(ctx)

	log.Info().Caller().Msg(fmt.Sprintf("=> worker %d inserted %d data", workerIndex, counter))
}

func (s Seed) AccessSeed() {
	uow, err := service.NewUnitOfWork(s.pool)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create unit of work")
	}

	ctx := context.Background()
	tx, err := uow.Begin(ctx)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to begin transaction")
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	endpointDatas := readYAML("endpoints.yml")
	var endpointYAML EndpointYAML
	err = yaml.Unmarshal(endpointDatas, &endpointYAML)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to unmarshal endpoint data")
	}

	roleDatas := readYAML("roles.yml")
	var roleYAML []RoleYAML
	err = yaml.Unmarshal(roleDatas, &roleYAML)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to unmarshal role data")
	}

	cachedEndpoint := make(map[string]domain.Endpoint)

	for _, endpoint := range endpointYAML.Endpoints {
		endpointData := domain.NewEndpoint(endpoint.Name, endpoint.Path, endpoint.Method)
		cachedEndpoint[endpoint.Name] = endpointData
	}

	for _, role := range roleYAML {
		roleData := domain.NewRole(role.Name)
		for _, endpoint := range role.Endpoints {
			if val, ok := cachedEndpoint[endpoint.Name]; ok {
				roleData.Endpoints.Add(val)
			}
		}
		log.Info().Caller().Msg(fmt.Sprintf("=> inserting role %s with endpoints size %d", roleData.Name, len(roleData.Endpoints)))
		roleErr := uow.Role.Save(ctx, tx, roleData)
		if roleErr != nil {
			log.Error().Caller().Err(roleErr).Msg("Failed to insert role")
		}
	}

	tx.Commit(ctx)
}

func readYAML(filename string) []byte {
	filePath := filepath.Join("data", filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg(fmt.Sprintf("Failed to read role data %s", filename))
	}
	return data
}
