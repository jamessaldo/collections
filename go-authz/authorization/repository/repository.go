package repository

import "authorization/infrastructure/persistence"

var (
	User       UserRepository
	Role       RoleRepository
	Team       TeamRepository
	Membership MembershipRepository
	Invitation InvitationRepository
)

func CreateRepositories() {
	User = NewUserRepository(persistence.Pool)
	Role = NewRoleRepository(persistence.Pool)
	Team = NewTeamRepository(persistence.Pool)
	Membership = NewMembershipRepository(persistence.Pool)
	Invitation = NewInvitationRepository(persistence.Pool)
}
