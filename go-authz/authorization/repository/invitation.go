package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

type invitationRepository struct {
	pool *pgxpool.Pool
}

type InvitationRepository interface {
	Add(*domain.Invitation, pgx.Tx) (*domain.Invitation, error)
	AddBatch([]domain.Invitation) error
	Update(*domain.Invitation, pgx.Tx) (*domain.Invitation, error)
	Get(ulid.ULID) (*domain.Invitation, error)
	List(opts *domain.InvitationOptions) ([]*domain.Invitation, error)
	Delete(ulid.ULID, pgx.Tx) error
}

// invitationRepository implements the InvitationRepository interface
func NewInvitationRepository(pool *pgxpool.Pool) InvitationRepository {
	return &invitationRepository{pool: pool}
}

func (repo *invitationRepository) Add(invitation *domain.Invitation, tx pgx.Tx) (*domain.Invitation, error) {
	query := `
		INSERT INTO invitations (id, email, expires_at, status, team_id, role_id, sender_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	_, err := tx.Exec(context.Background(), query,
		invitation.ID, invitation.Email, invitation.ExpiresAt, invitation.Status,
		invitation.TeamID, invitation.RoleID, invitation.SenderID, invitation.IsActive,
	)
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

// add batch pgx
func (repo *invitationRepository) AddBatch(invitations []domain.Invitation) error {
	query := `
		INSERT INTO invitations (id, email, expires_at, status, team_id, role_id, sender_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`

	tx, err := repo.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, invitation := range invitations {
		_, err := tx.Exec(context.Background(), query,
			invitation.ID, invitation.Email, invitation.ExpiresAt, invitation.Status,
			invitation.TeamID, invitation.RoleID, invitation.SenderID, invitation.IsActive,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (repo *invitationRepository) Update(invitation *domain.Invitation, tx pgx.Tx) (*domain.Invitation, error) {
	query := `
		UPDATE invitations
		SET email = $1, expires_at = $2, status = $3, team_id = $4, role_id = $5, sender_id = $6, is_active = $7
		WHERE id = $8
	`
	_, err := tx.Exec(context.Background(), query,
		invitation.Email, invitation.ExpiresAt, invitation.Status,
		invitation.TeamID, invitation.RoleID, invitation.SenderID, invitation.IsActive,
		invitation.ID,
	)
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

func (repo *invitationRepository) Get(id ulid.ULID) (*domain.Invitation, error) {
	var invitation domain.Invitation
	query := `
		SELECT i.id, i.email, i.expires_at, i.status, i.team_id, i.role_id, i.sender_id, i.is_active,
		FROM invitations i
		WHERE i.id = $1
	`
	err := repo.pool.QueryRow(context.Background(), query, id).
		Scan(&invitation.ID, &invitation.Email, &invitation.ExpiresAt, &invitation.Status,
			&invitation.TeamID, &invitation.RoleID, &invitation.SenderID, &invitation.IsActive,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}
	return &invitation, nil
}

func (repo *invitationRepository) List(opts *domain.InvitationOptions) ([]*domain.Invitation, error) {
	query := `
		SELECT id, email, expires_at, status, team_id, role_id, sender_id, is_active
		FROM invitations
		WHERE 1=1
	`
	var args []interface{}

	if len(opts.Statuses) > 0 {
		query += " AND status = ANY($1)"
		args = append(args, opts.Statuses)
	}
	if opts.Email != "" {
		query += " AND email = $2"
		args = append(args, opts.Email)
	}
	if opts.TeamID != uuid.Nil {
		query += " AND team_id = $3"
		args = append(args, opts.TeamID)
	}
	if !opts.ExpiresAt.IsZero() {
		query += " AND expires_at <= $4"
		args = append(args, opts.ExpiresAt)
	}
	if opts.RoleID.String() != "" {
		query += " AND role_id = $5"
		args = append(args, opts.RoleID)
	}
	if opts.Limit > 0 {
		query += " LIMIT $6"
		args = append(args, opts.Limit)
	}

	rows, err := repo.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []*domain.Invitation
	for rows.Next() {
		var invitation domain.Invitation
		err := rows.Scan(&invitation.ID, &invitation.Email, &invitation.ExpiresAt, &invitation.Status,
			&invitation.TeamID, &invitation.RoleID, &invitation.SenderID, &invitation.IsActive)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, &invitation)
	}
	return invitations, nil
}

func (repo *invitationRepository) Delete(id ulid.ULID, tx pgx.Tx) error {
	query := `
		DELETE FROM invitations
		WHERE id = $1
	`
	_, err := tx.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	return nil
}
