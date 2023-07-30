package service

// import (
// 	"authorization/repository"
// 	"context"

// 	"github.com/jackc/pgx/v5"         // Import the pgx package for working with pgx transactions
// 	"github.com/jackc/pgx/v5/pgxpool" // Import pgxpool for connection pooling
// )

// type UnitOfWork struct {
// 	pool       *pgxpool.Pool // Use pgxpool for connection pooling
// 	User       repository.UserRepository
// 	Role       repository.RoleRepository
// 	Team       repository.TeamRepository
// 	Membership repository.MembershipRepository
// 	Invitation repository.InvitationRepository
// }

// func NewUnitOfWork(pool *pgxpool.Pool) (*UnitOfWork, error) {
// 	return &UnitOfWork{
// 		pool:       pool,
// 		User:       repository.NewUserRepository(pool),
// 		Role:       repository.NewRoleRepository(pool),
// 		Team:       repository.NewTeamRepository(pool),
// 		Membership: repository.NewMembershipRepository(pool),
// 		Invitation: repository.NewInvitationRepository(pool),
// 	}, nil
// }

// func (u *UnitOfWork) GetDB() *pgxpool.Pool {
// 	return u.pool
// }

// func (u *UnitOfWork) Begin(ctx context.Context) (pgx.Tx, error) {
// 	tx, err := u.pool.Begin(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tx, nil
// }
