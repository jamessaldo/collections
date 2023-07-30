package seeder

import (
	"authorization/domain"
	"authorization/repository"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

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
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to begin transaction")
	}

	userRepo := repository.NewUserRepository(pool)
	// for _, user := range user_list {
	// 	_, userErr := userRepo.Add(user, tx)
	// 	if userErr != nil {
	// 		log.Error().Caller().Err(err).Msg("Failed to insert users")
	// 	}
	// }
	userErr := userRepo.AddBatch(ctx, user_list, tx)
	if userErr != nil {
		log.Error().Caller().Err(err).Msg("Failed to insert users")
	}
	tx.Commit(ctx)

	log.Info().Caller().Msg(fmt.Sprintf("=> worker %d inserted %d data", workerIndex, counter))
}

func _(pool *pgxpool.Pool, user_data domain.User, counter, workerIndex int) {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to begin transaction")
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	userRepo := repository.NewUserRepository(pool)

	_, userErr := userRepo.Add(ctx, user_data, tx)
	if userErr != nil {
		log.Error().Caller().Err(err).Msg("Failed to insert user")
		tx.Rollback(ctx)
	}

	tx.Commit(ctx)

	log.Info().Caller().Msg(fmt.Sprintf("=> worker %d inserted %d data", workerIndex, counter))
}
