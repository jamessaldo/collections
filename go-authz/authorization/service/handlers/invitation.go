package handlers

import (
	"errors"
	"fmt"
	"time"

	"auth/controller/exception"
	"auth/domain/command"
	"auth/domain/model"
	"auth/infrastructure/worker"
	"auth/service"

	uuid "github.com/satori/go.uuid"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	} else if team.IsPersonal {
		return exception.NewForbiddenException(fmt.Sprintf("you can't invite a member to personal team with ID %s", cmd.TeamID))
	}

	membershipOpts := &model.MembershipOptions{
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.NewNotFoundException(err.Error())
			}
			return err
		}

		inviteesOpts := &model.InvitationOptions{
			Email:    invitee.Email,
			TeamID:   cmd.TeamID,
			RoleID:   role.ID,
			Statuses: []model.InvitationStatus{model.InvitationStatusPending, model.InvitationStatusSent},
		}

		activeInvitees, err := uow.Invitation.List(inviteesOpts)
		if err != nil {
			return err
		}

		if len(activeInvitees) > 0 {
			continue
		}

		// create invitation
		invitation := model.Invitation{
			ID:        ksuid.New().String(),
			Email:     invitee.Email,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
			Status:    model.InvitationStatusPending,
			TeamID:    cmd.TeamID,
			RoleID:    role.ID,
			SenderID:  cmd.Sender.ID,
		}

		// add invitation
		_, err = uow.Invitation.Add(&invitation)
		if err != nil {
			return err
		}

		emailPayload := &worker.Payload{
			TemplateName: "invitation-message.html",
			To:           invitation.Email,
			Subject:      fmt.Sprintf("Invitation to join %s team", team.Name),
			Data: map[string]interface{}{
				"SenderName":     cmd.Sender.FullName(),
				"TeamName":       team.Name,
				"EmailTo":        invitation.Email,
				"InvitationLink": fmt.Sprintf("http://localhost:3000/invitation/%s", invitation.ID),
				"InvitationID":   invitation.ID,
			},
		}

		// send email
		errSendMail := mailer.SendEmail(emailPayload)
		if errSendMail != nil {
			return errSendMail
		}
	}
	tx.Commit()
	return nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	} else if team.IsPersonal {
		return exception.NewForbiddenException(fmt.Sprintf("you can't invite a member to personal team with ID %s", cmd.TeamID))
	}

	// get invitation
	invitation, err := uow.Invitation.Get(cmd.InvitationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	// check if invitation is not expired
	if invitation.ExpiresAt.After(time.Now()) || invitation.Status != model.InvitationStatusExpired {
		return exception.NewBadRequestException("invitation is not expired")
	}

	invitation.Status = model.InvitationStatusPending
	invitation.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)

	emailPayload := &worker.Payload{
		TemplateName: "invitation-message.html",
		To:           invitation.Email,
		Subject:      fmt.Sprintf("Invitation to join %s team", team.Name),
		Data: map[string]interface{}{
			"SenderName":     cmd.Sender.FullName(),
			"TeamName":       team.Name,
			"EmailTo":        invitation.Email,
			"InvitationLink": fmt.Sprintf("http://localhost:3000/invitation/%s", invitation.ID),
			"InvitationID":   invitation.ID,
		},
	}

	// send email
	errSendMail := mailer.SendEmail(emailPayload)
	if errSendMail != nil {
		return errSendMail
	}

	return nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	if invitation.Status == model.InvitationStatusAccepted {
		return exception.NewBadRequestException("invitation already accepted")
	}

	if invitation.Status == model.InvitationStatusDeclined {
		return exception.NewBadRequestException("invitation already declined")
	}

	// delete invitation
	err = uow.Invitation.Delete(invitation.ID)
	if err != nil {
		return err
	}

	return nil
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
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
	invitation.Status = model.InvitationStatus(cmd.Status)
	invitation.IsActive = false
	invitation, err = uow.Invitation.Update(invitation)
	if err != nil {
		return err
	}

	// add team member
	if cmd.Status == string(model.InvitationStatusAccepted) {
		role, err := uow.Role.Get(invitation.Role.Name)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return exception.NewNotFoundException(err.Error())
			}
			return err
		}
		membership := &model.Membership{
			ID:     uuid.NewV4(),
			TeamID: invitation.TeamID,
			UserID: cmd.User.ID,
			RoleID: role.ID,
		}

		_, err = uow.Membership.Add(membership)
		if err != nil {
			return err
		}
	}

	return nil
}
