package repository

import (
	"authorization/controller/exception"
	"authorization/domain"

	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type membershipRepository struct {
	pool *pgxpool.Pool // Use pgxpool.Pool for connection pooling
}

type MembershipRepository interface {
	Add(*domain.Membership, pgx.Tx) (*domain.Membership, error)
	AddBatch([]domain.Membership) error
	Update(*domain.Membership, pgx.Tx) (*domain.Membership, error)
	Get(uuid.UUID) (*domain.Membership, error)
	List(opts *domain.MembershipOptions) ([]domain.Membership, error)
	Delete(uuid.UUID, pgx.Tx) error
	Count(opts *domain.MembershipOptions) (int64, error)
}

// membershipRepository implements the MembershipRepository interface
func NewMembershipRepository(pool *pgxpool.Pool) MembershipRepository {
	return &membershipRepository{pool: pool}
}

func (repo *membershipRepository) Add(membership *domain.Membership, tx pgx.Tx) (*domain.Membership, error) {
	query := `
		INSERT INTO memberships (id, team_id, user_id, role_id, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := tx.QueryRow(
		context.Background(),
		query,
		membership.ID,
		membership.TeamID,
		membership.UserID,
		membership.RoleID,
		membership.LastActiveAt,
		membership.CreatedAt,
		membership.UpdatedAt,
	).Scan(&membership.ID)

	if err != nil {
		return nil, err
	}

	return membership, nil
}

func (repo *membershipRepository) AddBatch(memberships []domain.Membership) error {
	query := `
		INSERT INTO memberships (id, team_id, user_id, role_id, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`

	tx, err := repo.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, membership := range memberships {
		if _, err := tx.Exec(
			context.Background(),
			query,
			membership.ID,
			membership.TeamID,
			membership.UserID,
			membership.RoleID,
			membership.LastActiveAt,
			membership.CreatedAt,
			membership.UpdatedAt,
		); err != nil {
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func (repo *membershipRepository) Update(membership *domain.Membership, tx pgx.Tx) (*domain.Membership, error) {
	query := `
		UPDATE memberships
		SET team_id = $2, user_id = $3, role_id = $4, last_active_at = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := tx.Exec(
		context.Background(),
		query,
		membership.ID,
		membership.TeamID,
		membership.UserID,
		membership.RoleID,
		membership.LastActiveAt,
		membership.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return membership, nil
}

func (repo *membershipRepository) Get(id uuid.UUID) (*domain.Membership, error) {
	query := `
		SELECT id, team_id, user_id, role_id, last_active_at, created_at, updated_at
		FROM memberships
		WHERE id = $1
	`

	var membership domain.Membership

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&membership.ID,
		&membership.TeamID,
		&membership.UserID,
		&membership.RoleID,
		&membership.LastActiveAt,
		&membership.CreatedAt,
		&membership.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.NewNotFoundException("membership not found")
		}
		return nil, err
	}

	return &membership, nil
}

func (repo *membershipRepository) List(opts *domain.MembershipOptions) ([]domain.Membership, error) {
	query := `
		SELECT id, team_id, user_id, role_id, last_active_at, created_at, updated_at
		FROM memberships
		WHERE team_id = $1
	`

	var memberships []domain.Membership

	rows, err := repo.pool.Query(
		context.Background(),
		query,
		opts.TeamID,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var membership domain.Membership

		err := rows.Scan(
			&membership.ID,
			&membership.TeamID,
			&membership.UserID,
			&membership.RoleID,
			&membership.LastActiveAt,
			&membership.CreatedAt,
			&membership.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, nil
}

func (repo *membershipRepository) Delete(id uuid.UUID, tx pgx.Tx) error {
	query := `
		DELETE FROM memberships WHERE id = $1
	`

	_, err := tx.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *membershipRepository) Count(opts *domain.MembershipOptions) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM memberships
		WHERE team_id = $1
	`

	var count int64

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		opts.TeamID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
