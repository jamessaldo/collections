package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"fmt"

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
	Add(context.Context, domain.Membership, pgx.Tx) (domain.Membership, error)
	AddBatch(context.Context, []domain.Membership) error
	Update(context.Context, domain.Membership, pgx.Tx) (domain.Membership, error)
	Get(context.Context, uuid.UUID) (domain.Membership, error)
	List(context.Context, domain.MembershipOptions) ([]domain.Membership, error)
	Delete(context.Context, uuid.UUID, pgx.Tx) error
	Count(context.Context, domain.MembershipOptions) (int64, error)
}

// membershipRepository implements the MembershipRepository interface
func NewMembershipRepository(pool *pgxpool.Pool) MembershipRepository {
	return &membershipRepository{pool: pool}
}

func (repo *membershipRepository) Add(ctx context.Context, membership domain.Membership, tx pgx.Tx) (domain.Membership, error) {
	query := `
		INSERT INTO memberships (id, team_id, user_id, role_id, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := tx.QueryRow(
		ctx,
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
		return domain.Membership{}, err
	}

	return membership, nil
}

func (repo *membershipRepository) AddBatch(ctx context.Context, memberships []domain.Membership) error {
	query := `
		INSERT INTO memberships (id, team_id, user_id, role_id, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO NOTHING
	`

	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, membership := range memberships {
		if _, err := tx.Exec(
			ctx,
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

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (repo *membershipRepository) Update(ctx context.Context, membership domain.Membership, tx pgx.Tx) (domain.Membership, error) {
	query := `
		UPDATE memberships
		SET team_id = $2, user_id = $3, role_id = $4, last_active_at = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := tx.Exec(
		ctx,
		query,
		membership.ID,
		membership.TeamID,
		membership.UserID,
		membership.RoleID,
		membership.LastActiveAt,
		membership.UpdatedAt,
	)

	if err != nil {
		return domain.Membership{}, err
	}

	return membership, nil
}

func (repo *membershipRepository) Get(ctx context.Context, id uuid.UUID) (domain.Membership, error) {
	query := `
		SELECT id, team_id, user_id, role_id, last_active_at, created_at, updated_at
		FROM memberships
		WHERE id = $1
	`

	var membership domain.Membership

	err := repo.pool.QueryRow(
		ctx,
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
			return domain.Membership{}, exception.NewNotFoundException("membership not found")
		}
		return domain.Membership{}, err
	}

	return membership, nil
}

func (repo *membershipRepository) List(ctx context.Context, opts domain.MembershipOptions) ([]domain.Membership, error) {
	query := `
		SELECT id, team_id, user_id, role_id, last_active_at, created_at, updated_at
		FROM memberships
		WHERE
	`

	argsNumber := 1
	args := make([]interface{}, 0)

	if opts.TeamID != uuid.Nil {
		query += fmt.Sprintf("team_id = $%d", argsNumber)
		argsNumber++
		args = append(args, opts.TeamID)
	}

	if opts.UserID != uuid.Nil {
		query += fmt.Sprintf("user_id = $%d", argsNumber)
		argsNumber++
		args = append(args, opts.UserID)
	}

	var memberships []domain.Membership

	rows, err := repo.pool.Query(
		ctx,
		query,
		args...,
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

func (repo *membershipRepository) Delete(ctx context.Context, id uuid.UUID, tx pgx.Tx) error {
	query := `
		DELETE FROM memberships WHERE id = $1
	`

	_, err := tx.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *membershipRepository) Count(ctx context.Context, opts domain.MembershipOptions) (int64, error) {
	query := `
		SELECT COUNT(id)
		FROM memberships
		WHERE 
	`

	argsNumber := 1
	args := make([]interface{}, 0)

	if opts.TeamID != uuid.Nil {
		query += fmt.Sprintf("team_id = $%d", argsNumber)
		argsNumber++
		args = append(args, opts.TeamID)
	}

	if opts.UserID != uuid.Nil {
		query += fmt.Sprintf("user_id = $%d", argsNumber)
		argsNumber++
		args = append(args, opts.UserID)
	}

	var count int64

	err := repo.pool.QueryRow(
		ctx,
		query,
		args...,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
