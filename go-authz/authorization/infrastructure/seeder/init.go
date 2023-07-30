package seeder

import (
	"authorization/util"
	"fmt"
	"math"
	"reflect"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	pool *pgxpool.Pool
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
