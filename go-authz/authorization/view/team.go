package view

import (
	"authorization/controller/exception"
	"authorization/domain"
	"authorization/domain/dto"
	"authorization/repository"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	uuid "github.com/satori/go.uuid"
)

func Team(ctx context.Context, id uuid.UUID, user domain.User) (*dto.TeamRetrievalSchema, error) {
	team, err := repository.Team.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}
	opts := domain.MembershipOptions{
		TeamID:       id,
		IsSelectUser: true,
		IsSelectRole: true,
	}
	memberships, err := repository.Membership.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	totalMemberships, err := repository.Membership.Count(ctx, opts)
	if err != nil {
		return nil, err
	}
	var lastActiveAt time.Time
	var membershipsList []dto.MembershipRetrievalSchema

	for _, membership := range memberships {
		if membership.User.ID == user.ID {
			lastActiveAt = membership.LastActiveAt
		}
		membershipsList = append(membershipsList, membership.Parse())
	}

	return &dto.TeamRetrievalSchema{
		ID:           team.ID,
		Name:         team.Name,
		Description:  team.Description,
		AvatarURL:    team.AvatarURL,
		IsPersonal:   team.IsPersonal,
		Creator:      team.Creator.PublicUser(),
		LastActiveAt: lastActiveAt,
		NumOfMembers: totalMemberships,
		Memberships:  membershipsList,
		CreatedAt:    team.CreatedAt,
		UpdatedAt:    team.UpdatedAt,
	}, nil
}

func Teams(ctx context.Context, user domain.User, name string, page, pageSize int) (dto.Pagination, error) {
	var teams []interface{}
	var totalMemberships int64

	membershipOpts := domain.MembershipOptions{
		UserID:       user.ID,
		IsSelectTeam: true,
		Limit:        pageSize,
		Skip:         page,
		Name:         name,
	}
	memberships, err := repository.Membership.List(ctx, membershipOpts)
	if err != nil {
		return dto.Pagination{}, err
	}

	totalMemberships, err = repository.Membership.Count(ctx, membershipOpts)
	if err != nil {
		return dto.Pagination{}, err
	}
	for _, membership := range memberships {
		opts := domain.MembershipOptions{
			TeamID:       membership.TeamID,
			Limit:        3,
			IsSelectUser: true,
			IsSelectRole: true,
		}
		members, err := repository.Membership.List(ctx, opts)
		if err != nil {
			return dto.Pagination{}, err
		}
		totalMembers, err := repository.Membership.Count(ctx, opts)
		if err != nil {
			return dto.Pagination{}, err
		}

		var memberList []dto.MembershipRetrievalSchema
		for _, member := range members {
			data := dto.MembershipRetrievalSchema{
				ID:   member.ID,
				Role: string(member.Role.Name),
				User: member.User.PublicUser(),
			}
			memberList = append(memberList, data)
		}
		teams = append(teams, dto.TeamRetrievalSchema{
			ID:           membership.Team.ID,
			Name:         membership.Team.Name,
			Description:  membership.Team.Description,
			AvatarURL:    membership.Team.AvatarURL,
			IsPersonal:   membership.Team.IsPersonal,
			Creator:      membership.Team.Creator.PublicUser(),
			LastActiveAt: membership.LastActiveAt,
			NumOfMembers: totalMembers,
			Memberships:  memberList,
			CreatedAt:    membership.Team.CreatedAt,
			UpdatedAt:    membership.Team.UpdatedAt,
		})
	}
	return dto.Paginate(page, pageSize, totalMemberships, teams), nil
}
