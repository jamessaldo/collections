package view

import (
	"authorization/controller/exception"
	"authorization/domain/dto"
	"authorization/domain/model"
	"authorization/service"
	"errors"

	"gorm.io/gorm"
)

func Invitation(id string, uow *service.UnitOfWork) (*dto.InvitationRetreivalSchema, error) {
	invitation, err := uow.Invitation.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	team, err := uow.Team.Get(invitation.TeamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	membershipOpts := &model.MembershipOptions{
		TeamID:       team.ID,
		IsSelectUser: true,
		IsSelectRole: true,
		Limit:        3,
	}

	memberships, err := uow.Membership.List(membershipOpts)
	if err != nil {
		return nil, err
	}

	totalMemberships, err := uow.Membership.Count(membershipOpts)
	if err != nil {
		return nil, err
	}

	var membershipsList []*dto.MembershipRetrievalSchema

	for _, membership := range memberships {
		membershipsList = append(membershipsList, membership.Parse())
	}

	return &dto.InvitationRetreivalSchema{
		ID:         invitation.ID,
		Email:      invitation.Email,
		ExpiresAt:  invitation.ExpiresAt,
		Status:     string(invitation.Status),
		Role:       string(invitation.Role.Name),
		SenderName: invitation.Sender.FullName(),
		Team: dto.TeamRetrievalSchema{
			ID:           team.ID,
			Name:         team.Name,
			Description:  team.Description,
			AvatarURL:    team.AvatarURL,
			NumOfMembers: totalMemberships,
			Memberships:  membershipsList,
			Creator:      team.Creator.PublicUser(),
		},
	}, nil
}
