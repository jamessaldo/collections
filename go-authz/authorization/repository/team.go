package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"authorization/util"

	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type teamRepository struct {
	pool *pgxpool.Pool
}

// teamRepository implements the TeamRepository interface
type TeamRepository interface {
	Add(context.Context, domain.Team, pgx.Tx) (domain.Team, error)
	Update(context.Context, domain.Team, pgx.Tx) (domain.Team, error)
	Get(context.Context, uuid.UUID) (domain.Team, error)
}

func NewTeamRepository(pool *pgxpool.Pool) TeamRepository {
	return &teamRepository{pool: pool}
}

func (repo *teamRepository) Add(ctx context.Context, team domain.Team, tx pgx.Tx) (domain.Team, error) {
	query := `
		INSERT INTO teams (id, name, description, is_personal, avatar_url, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	_, err := tx.Exec(
		ctx,
		query,
		team.ID,
		team.Name,
		team.Description,
		team.IsPersonal,
		team.AvatarURL,
		team.CreatorID,
		team.CreatedAt,
		team.UpdatedAt,
	)

	if err != nil {
		return domain.Team{}, err
	}

	query = `
		INSERT INTO memberships (id, team_id, user_id, role_id, last_active_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	_, err = tx.Exec(
		ctx,
		query,
		team.Memberships[0].ID,
		team.ID,
		team.CreatorID,
		team.Memberships[0].RoleID,
		util.GetTimestampUTC(),
		util.GetTimestampUTC(),
		util.GetTimestampUTC(),
	)

	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

func (repo *teamRepository) Update(ctx context.Context, team domain.Team, tx pgx.Tx) (domain.Team, error) {

	query := `
		UPDATE teams
		SET name = $2, description = $3, is_personal = $4, avatar_url = $5, creator_id = $6, updated_at = $7
		WHERE id = $1	
	`

	_, err := tx.Exec(
		ctx,
		query,
		team.ID,
		team.Name,
		team.Description,
		team.IsPersonal,
		team.AvatarURL,
		team.CreatorID,
		team.UpdatedAt,
	)

	if err != nil {
		return domain.Team{}, err
	}

	for _, membership := range team.Memberships {
		q := `
			INSERT INTO memberships (id, team_id, user_id, role_id, last_active_at, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (id) DO NOTHING
		`

		_, err = tx.Exec(
			ctx,
			q,
			membership.ID,
			team.ID,
			team.CreatorID,
			membership.RoleID,
			util.GetTimestampUTC(),
			util.GetTimestampUTC(),
			util.GetTimestampUTC(),
		)

		if err != nil {
			return domain.Team{}, err
		}
	}

	return team, nil
}

func (repo *teamRepository) Get(ctx context.Context, id uuid.UUID) (domain.Team, error) {
	query := `
		SELECT id, name, description, is_personal, avatar_url, creator_id, created_at, updated_at
		FROM teams
		WHERE id = $1
	`

	var team domain.Team

	err := repo.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&team.ID,
		&team.Name,
		&team.Description,
		&team.IsPersonal,
		&team.AvatarURL,
		&team.CreatorID,
		&team.CreatedAt,
		&team.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Team{}, exception.NewNotFoundException("team not found")
		}
		return domain.Team{}, err
	}

	return team, nil
}
