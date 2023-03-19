package view

import (
	"auth/controller/exception"
	"auth/domain/dto"
	"auth/domain/model"
	"auth/service"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func Team(id uuid.UUID, user *model.User, uow *service.UnitOfWork) (*dto.TeamRetrievalSchema, error) {
	team, err := uow.Team.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}
	opts := &model.MembershipOptions{
		TeamID:       id,
		IsSelectUser: true,
		IsSelectRole: true,
	}
	memberships, err := uow.Membership.List(opts)
	if err != nil {
		return nil, err
	}
	totalMemberships, err := uow.Membership.Count(opts)
	if err != nil {
		return nil, err
	}
	var lastActiveAt *time.Time
	var membershipsList []*dto.MembershipRetrievalSchema

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

func Teams(uow *service.UnitOfWork, user *model.User, name string, page, pageSize int) (dto.Pagination, error) {
	var teams []*dto.TeamRetrievalSchema
	var totalMemberships int64

	membershipOpts := &model.MembershipOptions{
		UserID:       user.ID,
		IsSelectTeam: true,
		Limit:        pageSize,
		Skip:         page,
		Name:         name,
	}
	memberships, err := uow.Membership.List(membershipOpts)
	if err != nil {
		return dto.Pagination{}, err
	}

	totalMemberships, err = uow.Membership.Count(membershipOpts)
	if err != nil {
		return dto.Pagination{}, err
	}
	for _, membership := range memberships {
		opts := &model.MembershipOptions{
			TeamID:       membership.TeamID,
			Limit:        3,
			IsSelectUser: true,
			IsSelectRole: true,
		}
		members, err := uow.Membership.List(opts)
		if err != nil {
			return dto.Pagination{}, err
		}
		totalMembers, err := uow.Membership.Count(opts)
		if err != nil {
			return dto.Pagination{}, err
		}

		var memberList []*dto.MembershipRetrievalSchema
		for _, member := range members {
			data := dto.MembershipRetrievalSchema{
				ID:   member.ID,
				Role: string(member.Role.Name),
				User: member.User.PublicUser(),
			}
			memberList = append(memberList, &data)
		}
		teams = append(teams, &dto.TeamRetrievalSchema{
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
