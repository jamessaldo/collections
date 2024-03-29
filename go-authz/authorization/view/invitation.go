package view

import (
	"authorization/controller/exception"
	"authorization/domain"
	"authorization/domain/dto"
	"authorization/repository"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

func Invitation(ctx context.Context, id ulid.ULID) (*dto.InvitationRetreivalSchema, error) {
	invitation, err := repository.Invitation.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	team, err := repository.Team.Get(ctx, invitation.TeamID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	membershipOpts := domain.MembershipOptions{
		TeamID:       team.ID,
		IsSelectUser: true,
		IsSelectRole: true,
		Limit:        3,
	}

	memberships, err := repository.Membership.List(ctx, membershipOpts)
	if err != nil {
		return nil, err
	}

	totalMemberships, err := repository.Membership.Count(ctx, membershipOpts)
	if err != nil {
		return nil, err
	}

	var membershipsList []dto.MembershipRetrievalSchema

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
