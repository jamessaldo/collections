package domain

import (
	"authorization/controller/exception"
	"authorization/domain/dto"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

type Team struct {
	ID          uuid.UUID
	Name        string
	Description string
	IsPersonal  bool
	AvatarURL   string
	CreatorID   uuid.UUID
	Creator     *User
	Memberships []Membership
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Membership struct {
	ID     uuid.UUID
	TeamID uuid.UUID
	Team   *Team
	UserID uuid.UUID
	User   *User
	RoleID ulid.ULID
	Role   *Role

	LastActiveAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *Membership) Parse() *dto.MembershipRetrievalSchema {
	return &dto.MembershipRetrievalSchema{
		ID:   m.ID,
		User: m.User.PublicUser(),
		Role: string(m.Role.Name),
	}
}

type MembershipOptions struct {
	Limit        int
	Skip         int
	Name         string
	IsSelectTeam bool
	IsSelectUser bool
	IsSelectRole bool
	TeamID       uuid.UUID
	UserID       uuid.UUID
	RoleID       uuid.UUID
}

func (t *Team) Update(payload map[string]any) {
	if val, ok := payload["name"].(string); ok && val != "" {
		t.Name = val
	}

	if val, ok := payload["description"].(string); ok && val != "" {
		t.Description = val
	}

	if val, ok := payload["avatarURL"].(string); ok && val != "" {
		t.AvatarURL = val
	}
}

func (t *Team) AddMembership(teamID, userID uuid.UUID, roleID ulid.ULID) {
	membership := Membership{
		ID:     uuid.NewV4(),
		TeamID: teamID,
		UserID: userID,
		RoleID: roleID,
	}
	t.Memberships = append(t.Memberships, membership)
}

func (m *Membership) Validation(userID, teamID uuid.UUID, requestedRole RoleType) error {
	if m.UserID == userID {
		return exception.NewForbiddenException("You cannot change your role")
	}

	if m.TeamID != teamID {
		return exception.NewForbiddenException(fmt.Sprintf("Team with ID %s is not match with membership-team ID", teamID))
	}

	if m.Role.Name == Owner {
		return exception.NewForbiddenException("It's not allowed to change owner role")
	}

	if requestedRole != "" && requestedRole == Owner {
		return exception.NewForbiddenException("You cannot change role to owner")
	}

	return nil
}

func NewTeam(user User, roleID ulid.ULID, name, description string, isPersonal bool) *Team {
	teamID := uuid.NewV4()

	membership := Membership{
		ID:     uuid.NewV4(),
		TeamID: teamID,
		UserID: user.ID,
		RoleID: roleID,
	}

	team := &Team{
		ID:          teamID,
		Name:        name,
		Description: description,
		IsPersonal:  isPersonal,
		CreatorID:   user.ID,
		Memberships: []Membership{membership},
	}

	if isPersonal {
		team.Name = fmt.Sprintf("%s's Personal Team", user.FullName())
		team.Description = fmt.Sprintf("%s's Personal Team will contains your personal apps.", user.FullName())
	}

	return team
}
