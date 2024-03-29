package domain

import (
	"authorization/util"
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type RoleType string

var (
	Owner   RoleType = "owner"
	Admin   RoleType = "admin"
	Member  RoleType = "member"
	Finance RoleType = "finance"
)

type EndpointYAML struct {
	Endpoints []struct {
		Name   string `yaml:"name"`
		Path   string `yaml:"path"`
		Method string `yaml:"method"`
	} `yaml:"endpoints"`
}

type RoleYAML struct {
	Name      RoleType `yaml:"name"`
	Endpoints []struct {
		Name string `yaml:"name"`
	} `yaml:"endpoints"`
}

type Endpoint struct {
	Name   string
	Path   string
	Method string
}

func (e Endpoint) Equals(endpoint Endpoint) bool {
	return e.Name == endpoint.Name && e.Path == endpoint.Path && e.Method == endpoint.Method
}

type Endpoints []Endpoint

type Role struct {
	ID        ulid.ULID
	Name      RoleType
	Endpoints Endpoints
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (endpoints Endpoints) Find(name string) Endpoint {
	for _, e := range endpoints {
		if e.Name == name {
			return e
		}
	}
	return Endpoint{}
}

func (endpoints Endpoints) Contains(endpoint Endpoint) bool {
	for _, e := range endpoints {
		if e.Equals(endpoint) {
			return true
		}
	}
	return false
}

func (endpoints *Endpoints) Remove(endpoint Endpoint) {
	endpointsValue := *endpoints
	for i, e := range endpointsValue {
		if e.Name == endpoint.Name {
			*endpoints = append(endpointsValue[:i], endpointsValue[i+1:]...)
		}
	}
}

func (endpoints *Endpoints) Add(endpoint Endpoint) {
	*endpoints = append(*endpoints, endpoint)
}

func (endpoints Endpoints) ToJSON() ([]byte, error) {
	json, err := json.Marshal(endpoints)
	if err != nil {
		log.Error().Err(err).Msg("error marshaling endpoints")
		return nil, err
	}
	return json, nil
}

type Access struct {
	RoleName  RoleType
	IsAllowed bool
	Endpoint  Endpoint
}

func NewRole(name RoleType) Role {
	now := util.GetTimestampUTC()
	return Role{ID: ulid.Make(), Name: name, CreatedAt: now, UpdatedAt: now}
}

func NewEndpoint(name, path, method string) Endpoint {
	return Endpoint{Name: name, Path: path, Method: method}
}
