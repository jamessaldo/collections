package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type invitationRepository struct {
	pool *pgxpool.Pool
}

type InvitationRepository interface {
	Add(domain.Invitation, pgx.Tx) (domain.Invitation, error)
	AddBatch([]domain.Invitation) error
	Update(domain.Invitation, pgx.Tx) error
	Get(ulid.ULID) (domain.Invitation, error)
	List(opts domain.InvitationOptions) ([]domain.Invitation, error)
	Delete(ulid.ULID, pgx.Tx) error
}

// invitationRepository implements the InvitationRepository interface
func NewInvitationRepository(pool *pgxpool.Pool) InvitationRepository {
	return &invitationRepository{pool: pool}
}

func (repo *invitationRepository) Add(invitation domain.Invitation, tx pgx.Tx) (domain.Invitation, error) {
	query := `
		INSERT INTO invitations (id, email, expires_at, status, team_id, role_id, sender_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	_, err := tx.Exec(context.Background(), query,
		invitation.ID, invitation.Email, invitation.ExpiresAt, invitation.Status,
		invitation.TeamID, invitation.RoleID, invitation.SenderID, invitation.IsActive,
		invitation.CreatedAt, invitation.UpdatedAt,
	)
	if err != nil {
		return domain.Invitation{}, err
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

func (repo *invitationRepository) Update(invitation domain.Invitation, tx pgx.Tx) error {
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
		return err
	}

	return nil
}

func (repo *invitationRepository) Get(id ulid.ULID) (domain.Invitation, error) {
	var invitation domain.Invitation
	query := `
		SELECT i.id, i.email, i.expires_at, i.status, i.team_id, i.role_id, i.sender_id, i.is_active
		FROM invitations i
		WHERE i.id = $1
	`
	err := repo.pool.QueryRow(context.Background(), query, id).
		Scan(&invitation.ID, &invitation.Email, &invitation.ExpiresAt, &invitation.Status,
			&invitation.TeamID, &invitation.RoleID, &invitation.SenderID, &invitation.IsActive,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Invitation{}, exception.NewNotFoundException(err.Error())
		}
		return domain.Invitation{}, err
	}
	return invitation, nil
}

func (repo *invitationRepository) List(opts domain.InvitationOptions) ([]domain.Invitation, error) {
	query := `
		SELECT id, email, expires_at, status, team_id, role_id, sender_id, is_active
		FROM invitations
		WHERE status = ANY($1) AND email = $2 AND team_id = $3 AND role_id = $4	
	`
	rows, err := repo.pool.Query(context.Background(), query, opts.Statuses, opts.Email, opts.TeamID, opts.RoleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []domain.Invitation
	for rows.Next() {
		var invitation domain.Invitation
		err := rows.Scan(&invitation.ID, &invitation.Email, &invitation.ExpiresAt, &invitation.Status,
			&invitation.TeamID, &invitation.RoleID, &invitation.SenderID, &invitation.IsActive)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
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
