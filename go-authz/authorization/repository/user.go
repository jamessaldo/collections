package repository

import (
	"authorization/domain"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	uuid "github.com/satori/go.uuid"
)

type userRepository struct {
	pool *pgxpool.Pool // Use pgxpool.Pool for connection pooling
}

type UserRepository interface {
	Add(context.Context, domain.User, pgx.Tx) (domain.User, error)
	AddBatch(context.Context, domain.Users, pgx.Tx) error
	Update(context.Context, domain.User, pgx.Tx) (domain.User, error)
	Get(context.Context, uuid.UUID) (domain.User, error)
	List(context.Context, int, int) (domain.Users, error)
	GetByEmail(context.Context, string) (domain.User, error)
	GetByUsername(context.Context, string) (domain.User, error)
	Count(context.Context) (int64, error)
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (repo *userRepository) Add(ctx context.Context, user domain.User, tx pgx.Tx) (domain.User, error) {
	query := `INSERT INTO users (id, first_name, last_name, email, username, password, phone_number, avatar_url, is_active, verified, provider, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := tx.Exec(ctx, query, user.ID, user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.PhoneNumber, user.AvatarURL, user.IsActive, user.Verified, user.Provider, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (repo *userRepository) AddBatch(ctx context.Context, users domain.Users, tx pgx.Tx) error {
	// Create a slice of rows to hold the user data.
	var rows [][]interface{}
	for _, user := range users {
		rows = append(rows, []interface{}{
			user.ID,
			user.FirstName,
			user.LastName,
			user.Email,
			user.Username,
			user.Password,
			user.PhoneNumber,
			user.AvatarURL,
			user.IsActive,
			user.Verified,
			user.Provider,
			user.CreatedAt,
			user.UpdatedAt,
		})
	}

	// Perform the bulk insert using pgx.CopyFrom.
	_, err := tx.CopyFrom(ctx, pgx.Identifier{"users"}, []string{
		"id", "first_name", "last_name", "email", "username", "password", "phone_number", "avatar_url", "is_active", "verified", "provider", "created_at", "updated_at",
	}, pgx.CopyFromRows(rows))
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) Update(ctx context.Context, user domain.User, tx pgx.Tx) (domain.User, error) {
	query := "UPDATE users SET first_name = $1, last_name = $2, email = $3, username = $4, password = $5, phone_number = $6, avatar_url = $7, is_active = $8, verified = $9, provider = $10, updated_at = $11 WHERE id = $12"

	_, err := tx.Exec(ctx, query, user.FirstName, user.LastName, user.Email, user.Username, user.Password, user.PhoneNumber, user.AvatarURL, user.IsActive, user.Verified, user.Provider, user.UpdatedAt, user.ID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (repo *userRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	query := "SELECT id, first_name, last_name, email, username, password, phone_number, avatar_url, is_active, verified, provider, created_at, updated_at FROM users WHERE id = $1 AND is_active = true"

	var user domain.User
	err := repo.pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.PhoneNumber, &user.AvatarURL, &user.IsActive, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (repo *userRepository) List(ctx context.Context, page, pageSize int) (domain.Users, error) {
	query := "SELECT id, first_name, last_name, email, username, password, phone_number, avatar_url, is_active, verified, provider, created_at, updated_at FROM users WHERE is_active = true LIMIT $1 OFFSET $2"

	var users domain.Users
	rows, err := repo.pool.Query(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.PhoneNumber, &user.AvatarURL, &user.IsActive, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return domain.Users{}, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (repo *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := "SELECT id, first_name, last_name, email, username, password, phone_number, avatar_url, is_active, verified, provider, created_at, updated_at FROM users WHERE email = $1 AND is_active = true"

	var user domain.User
	err := repo.pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.PhoneNumber, &user.AvatarURL, &user.IsActive, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (repo *userRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	query := "SELECT id, first_name, last_name, email, username, password, phone_number, avatar_url, is_active, verified, provider, created_at, updated_at FROM users WHERE username = $1 AND is_active = true"

	var user domain.User
	err := repo.pool.QueryRow(ctx, query, username).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.Password, &user.PhoneNumber, &user.AvatarURL, &user.IsActive, &user.Verified, &user.Provider, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (repo *userRepository) Count(ctx context.Context) (int64, error) {
	query := "SELECT COUNT(*) FROM users WHERE is_active = true"

	var count int64
	err := repo.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
