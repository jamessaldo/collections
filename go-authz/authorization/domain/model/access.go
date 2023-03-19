package model

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type RoleType string

var (
	Owner   RoleType = "owner"
	Admin   RoleType = "admin"
	Member  RoleType = "member"
	Finance RoleType = "finance"
)

type Endpoint struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"size:100;not null;"`
	Path      string    `gorm:"size:100;not null;"`
	Method    string    `gorm:"size:100;not null;"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Role struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;"`
	Name      RoleType    `gorm:"size:100;not null;unique"`
	Endpoints []*Endpoint `gorm:"many2many:accesses;"`
	CreatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time   `gorm:"default:CURRENT_TIMESTAMP"`
}

type Access struct {
	ID         int64     `gorm:"primary_key;auto_increment"`
	RoleID     uuid.UUID `gorm:"type:uuid;not null;primaryKey;uniqueIndex:access_idx"`
	EndpointID uuid.UUID `gorm:"type:uuid;not null;primaryKey;uniqueIndex:access_idx"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func NewRole(name RoleType) *Role {
	return &Role{Name: name}
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
