package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type roleRepository struct {
	pool *pgxpool.Pool // Use pgxpool.Pool for connection pooling
}

// roleRepository implements the RoleRepository interface
type RoleRepository interface {
	Save(context.Context, pgx.Tx, *domain.Role) error
	Get(context.Context, domain.RoleType) (*domain.Role, error)
	GetAccess(context.Context, uuid.UUID, uuid.UUID, domain.Endpoint) (domain.Access, error)
}

func NewRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &roleRepository{pool: pool}
}

func (repo *roleRepository) Save(ctx context.Context, tx pgx.Tx, role *domain.Role) error {
	query := `
		INSERT INTO roles (id, name, endpoints, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			name = $2,
			endpoints = $3,
			updated_at = $5
	`

	endpoints, err := role.Endpoints.ToJSON()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		query,
		role.ID,
		role.Name,
		endpoints,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *roleRepository) Get(ctx context.Context, name domain.RoleType) (*domain.Role, error) {
	query := `
		SELECT id, name, endpoints, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	var role *domain.Role
	var endpointsJSON []byte

	err := repo.pool.QueryRow(
		ctx,
		query,
		name,
	).Scan(
		role.ID,
		role.Name,
		endpointsJSON,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.NewNotFoundException(fmt.Sprintf("Role with name %s does not exist", name))
		}
		return nil, err
	}

	err = json.Unmarshal(endpointsJSON, role.Endpoints)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal endpoints")
		return nil, err
	}

	return role, nil
}

func (repo *roleRepository) GetAccess(ctx context.Context, teamID, userID uuid.UUID, endpoint domain.Endpoint) (domain.Access, error) {
	query := `
		SELECT r.name, r.endpoints
		FROM memberships m
		LEFT JOIN roles r ON r.id = m.role_id
		WHERE m.team_id = $1 AND m.user_id = $2
	`

	var access domain.Access
	var endpointsJSON []byte
	var endpoints domain.Endpoints

	err := repo.pool.QueryRow(
		ctx,
		query,
		teamID,
		userID,
	).Scan(access.RoleName, endpointsJSON)

	if err != nil {
		return domain.Access{}, err
	}

	err = json.Unmarshal(endpointsJSON, endpoints)

	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal endpoints")
		return domain.Access{}, err
	}

	access.IsAllowed = endpoints.Contains(endpoint)
	access.Endpoint = endpoint
	return access, nil
}
