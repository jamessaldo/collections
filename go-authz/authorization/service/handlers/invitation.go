package handlers

import (
	"fmt"

	"authorization/controller/exception"
	"authorization/domain"
	"authorization/domain/command"
	"authorization/infrastructure/worker"
	"authorization/service"

	"gorm.io/gorm"
)

func InviteMemberWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.InviteMember); ok {
		return InviteMember(uow, mailer, c)
	}
	return fmt.Errorf("invalid command type, expected *command.InviteMember, got %T", cmd)
}

// create inviteMember function
func InviteMember(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd *command.InviteMember) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	// get team
	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		return err
	} else if team.IsPersonal {
		return exception.NewForbiddenException(fmt.Sprintf("you can't invite a member to personal team with ID %s", cmd.TeamID))
	}

	membershipOpts := &domain.MembershipOptions{
		TeamID:       cmd.TeamID,
		IsSelectUser: true,
		IsSelectRole: true,
	}
	memberships, err := uow.Membership.List(membershipOpts)
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

		role, err := uow.Role.Get(invitee.Role)
		if err != nil {
			return err
		}

		inviteesOpts := &domain.InvitationOptions{
			Email:    invitee.Email,
			TeamID:   cmd.TeamID,
			RoleID:   role.ID,
			Statuses: []domain.InvitationStatus{domain.InvitationStatusPending, domain.InvitationStatusSent},
		}

		activeInvitees, err := uow.Invitation.List(inviteesOpts)
		if err != nil {
			return err
		}

		if len(activeInvitees) > 0 {
			continue
		}

		invitation := domain.NewInvitation(invitee.Email, domain.InvitationStatusPending, cmd.TeamID, cmd.Sender.ID, role.ID)
		_, err = uow.Invitation.Add(invitation, tx)
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

		emailPayload := mailer.CreatePayload(worker.InvitationTemplate, invitation.Email, fmt.Sprintf("Invitation to join %s team", team.Name), data)

		// send email
		errSendMail := mailer.SendEmail(emailPayload)
		if errSendMail != nil {
			return errSendMail
		}
	}
	tx.Commit()
	return nil
}

func ResendInvitationWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.ResendInvitation); ok {
		return ResendInvitation(uow, mailer, c)
	}
	return fmt.Errorf("invalid command type, expected *command.ResendInvitation, got %T", cmd)
}

// create ResendInvitation function
func ResendInvitation(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd *command.ResendInvitation) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	// get team
	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		return err
	} else if team.IsPersonal {
		return exception.NewForbiddenException(fmt.Sprintf("you can't invite a member to personal team with ID %s", cmd.TeamID))
	}

	// get invitation
	invitation, err := uow.Invitation.Get(cmd.InvitationID)
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

	emailPayload := mailer.CreatePayload(worker.InvitationTemplate, invitation.Email, fmt.Sprintf("Invitation to join %s team", team.Name), data)

	// send email
	errSendMail := mailer.SendEmail(emailPayload)
	if errSendMail != nil {
		return errSendMail
	}

	return nil
}

func DeleteInvitationWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteInvitation); ok {
		return DeleteInvitation(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteInvitation, got %T", cmd)
}

// create DeleteInvitation function
func DeleteInvitation(uow *service.UnitOfWork, cmd *command.DeleteInvitation) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	// get invitation
	invitation, err := uow.Invitation.Get(cmd.InvitationID)
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
	err = uow.Invitation.Delete(invitation.ID, tx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateInvitationStatusWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateInvitationStatus); ok {
		return UpdateInvitationStatus(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateInvitationStatus, got %T", cmd)
}

// create UpdateInvitationStatus function
func UpdateInvitationStatus(uow *service.UnitOfWork, cmd *command.UpdateInvitationStatus) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	// get invitation
	invitation, err := uow.Invitation.Get(cmd.InvitationID)
	if err != nil {
		return err
	}

	team := &invitation.Team

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
	invitation, err = uow.Invitation.Update(invitation, tx)
	if err != nil {
		return err
	}

	// add team member
	if cmd.Status == string(domain.InvitationStatusAccepted) {
		role, err := uow.Role.Get(invitation.Role.Name)
		if err != nil {
			return err
		}

		team.AddMembership(invitation.TeamID, cmd.User.ID, role.ID)
		_, err = uow.Team.Update(team, tx)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}
