package seeder

import (
	"authorization/domain"
	"authorization/repository"
	"authorization/util"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

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

func (s Seed) AccessSeed() {
	ctx := context.Background()
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to begin transaction")
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	roleRepo := repository.NewRoleRepository(s.pool)

	endpointDatas := util.ReadYAML("endpoints.yml")
	var endpointYAML EndpointYAML
	err = yaml.Unmarshal(endpointDatas, &endpointYAML)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to unmarshal endpoint data")
	}

	roleDatas := util.ReadYAML("roles.yml")
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
		roleErr := roleRepo.Save(ctx, tx, roleData)
		if roleErr != nil {
			log.Error().Caller().Err(roleErr).Msg("Failed to insert role")
		}
	}

	tx.Commit(ctx)
}
