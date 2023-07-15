package model

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RoleType string

var (
	Owner   RoleType = "owner"
	Admin   RoleType = "admin"
	Member  RoleType = "member"
	Finance RoleType = "finance"
)

type Endpoint struct {
	ID        ulid.ULID `gorm:"type:bytea;primary_key"`
	Name      string    `gorm:"size:100;not null;unique;"`
	Path      string    `gorm:"size:100;not null;unique;"`
	Method    string    `gorm:"size:100;not null;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Role struct {
	ID        ulid.ULID   `gorm:"type:bytea;primary_key"`
	Name      RoleType    `gorm:"size:100;not null;unique"`
	Endpoints []*Endpoint `gorm:"many2many:accesses;"`
	CreatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
}

type Access struct {
	ID         int64     `gorm:"primary_key;auto_increment"`
	RoleID     ulid.ULID `gorm:"type:bytea;not null;primaryKey;uniqueIndex:access_idx"`
	EndpointID ulid.ULID `gorm:"type:bytea;not null;primaryKey;uniqueIndex:access_idx"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewRole(id ulid.ULID, name RoleType) *Role {
	return &Role{ID: id, Name: name}
}

func NewEndpoint(id ulid.ULID, name, path, method string) *Endpoint {
	return &Endpoint{ID: id, Name: name, Path: path, Method: method}
}

func (r *Role) AddEndpoints(endpoints ...*Endpoint) {
	r.Endpoints = append(r.Endpoints, endpoints...)
}

func (r *Role) RemoveEndpoint(endpoint *Endpoint) {
	for i, e := range r.Endpoints {
		if e.ID == endpoint.ID {
			r.Endpoints = append(r.Endpoints[:i], r.Endpoints[i+1:]...)
			break
		}
	}
}
