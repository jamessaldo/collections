package handlers

import (
	"authorization/controller/exception"
	"authorization/domain"
	"authorization/domain/command"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/worker"
	"authorization/repository"
	"context"
	"fmt"
)

// create inviteSendInvitationMember function
func SendInvitation(ctx context.Context, cmd *command.SendInvitation) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	// get team
	team, err := repository.Team.Get(ctx, cmd.TeamID)
	if err != nil {
		return err
	} else if team.IsPersonal {
		return exception.NewForbiddenException(fmt.Sprintf("you can't invite a member to personal team with ID %s", cmd.TeamID))
	}

	membershipOpts := domain.MembershipOptions{
		TeamID:       cmd.TeamID,
		IsSelectUser: true,
		IsSelectRole: true,
	}
	memberships, err := repository.Membership.List(ctx, membershipOpts)
	if err != nil {
		return err
	}

	emailSet := make(map[string]bool)

	for _, membership := range memberships {
		emailSet[membership.User.Email] = true
	}

	for _, invitee := range cmd.Invitees {
		if emailSet[invitee.Email] {
			continue
		}

		role, err := repository.Role.GetByName(ctx, invitee.Role)
		if err != nil {
			return err
		}

		inviteesOpts := domain.InvitationOptions{
			Email:    invitee.Email,
			TeamID:   cmd.TeamID,
			RoleID:   role.ID,
			Statuses: []domain.InvitationStatus{domain.InvitationStatusPending, domain.InvitationStatusSent},
		}

		activeInvitees, err := repository.Invitation.List(ctx, inviteesOpts)
		if err != nil {
			return err
		}

		if len(activeInvitees) > 0 {
			continue
		}

		invitation := domain.NewInvitation(invitee.Email, domain.InvitationStatusPending, cmd.TeamID, cmd.Sender.ID, role.ID)
		_, err = repository.Invitation.Add(ctx, invitation, tx)
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"SenderName":     cmd.Sender.FullName(),
			"TeamName":       team.Name,
			"EmailTo":        invitation.Email,
			"InvitationLink": fmt.Sprintf("http://localhost:3000/invitation/%s", invitation.ID),
			"InvitationID":   invitation.ID,
		}

		emailPayload := worker.Mailer.CreateEmailPayload(worker.InvitationTemplate, invitation.Email, fmt.Sprintf("Invitation to join %s team", team.Name), data)

		// send email
		errSendMail := worker.Mailer.SendEmail(emailPayload)
		if errSendMail != nil {
			return errSendMail
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// create ResendInvitation function
func ResendInvitation(ctx context.Context, cmd *command.ResendInvitation) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	// get team
	team, err := repository.Team.Get(ctx, cmd.TeamID)
	if err != nil {
		return err
	} else if team.IsPersonal {
		return exception.NewForbiddenException(fmt.Sprintf("you can't invite a member to personal team with ID %s", cmd.TeamID))
	}

	// get invitation
	invitation, err := repository.Invitation.Get(ctx, cmd.InvitationID)
	if err != nil {
		return err
	}

	err = invitation.ResendUpdate()
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"SenderName":     cmd.Sender.FullName(),
		"TeamName":       team.Name,
		"EmailTo":        invitation.Email,
		"InvitationLink": fmt.Sprintf("http://localhost:3000/invitation/%s", invitation.ID),
		"InvitationID":   invitation.ID,
	}

	emailPayload := worker.Mailer.CreateEmailPayload(worker.InvitationTemplate, invitation.Email, fmt.Sprintf("Invitation to join %s team", team.Name), data)

	// send email
	errSendMail := worker.Mailer.SendEmail(emailPayload)
	if errSendMail != nil {
		return errSendMail
	}

	err = repository.Invitation.Update(ctx, invitation, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// create DeleteInvitation function
func DeleteInvitation(ctx context.Context, cmd *command.DeleteInvitation) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	// get invitation
	invitation, err := repository.Invitation.Get(ctx, cmd.InvitationID)
	if err != nil {
		return err
	}

	if invitation.Status == domain.InvitationStatusAccepted {
		return exception.NewBadRequestException("invitation already accepted")
	}

	if invitation.Status == domain.InvitationStatusDeclined {
		return exception.NewBadRequestException("invitation already declined")
	}

	// delete invitation
	err = repository.Invitation.Delete(ctx, invitation.ID, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// create UpdateInvitationStatus function
func UpdateInvitationStatus(ctx context.Context, cmd *command.UpdateInvitationStatus) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	// get invitation
	invitation, err := repository.Invitation.Get(ctx, cmd.InvitationID)
	if err != nil {
		return err
	}

	team, err := repository.Team.Get(ctx, invitation.TeamID)
	if err != nil {
		return err
	}

	// check if user is team owner
	if invitation.Email != cmd.User.Email {
		return exception.NewForbiddenException(fmt.Sprintf("you are not invited to join this team with ID %s", invitation.TeamID))
	}

	if !invitation.IsActive {
		return exception.NewForbiddenException(fmt.Sprintf("invitation with ID %s is not active anymore", invitation.ID))
	}

	// update invitation
	invitation.Status = domain.InvitationStatus(cmd.Status)
	invitation.IsActive = false
	err = repository.Invitation.Update(ctx, invitation, tx)
	if err != nil {
		return err
	}

	// add team member
	if cmd.Status == string(domain.InvitationStatusAccepted) {
		role, err := repository.Role.Get(ctx, invitation.RoleID)
		if err != nil {
			return err
		}

		team.AddMembership(invitation.TeamID, cmd.User.ID, role.ID)
		_, err = repository.Team.Update(ctx, team, tx)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
