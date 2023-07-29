package integration

import (
	"authorization/domain"
	"authorization/domain/command"
	"authorization/view"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	uuid "github.com/satori/go.uuid"
)

func createTeam(cmd *command.CreateTeam, user domain.User) {
	ctx := context.Background()
	err := Bus.Handle(ctx, cmd)
	Ω(err).To(Succeed())

	team, err := view.Team(cmd.TeamID, user, Bus.UoW)
	Ω(err).To(Succeed())
	Ω(team.Name).To(Equal(cmd.Name))
	Ω(team.Description).To(Equal(cmd.Description))
	Ω(team.IsPersonal).To(BeFalse())
}

var _ = Describe("Team Testing", Ordered, func() {
	format.MaxLength = 0
	var (
		john domain.User
		jane domain.User
		cmdA *command.CreateTeam
		cmdB *command.CreateTeam
	)

	BeforeEach(func() {
		uow := Bus.UoW

		john = domain.NewUser("John", "Doe", "johndoe@example.com", "", "Google", true)
		err := createUser(john, uow)
		Ω(err).To(Succeed())

		jane = domain.NewUser("Jane", "Doe", "janedoe@example.com", "", "Google", true)
		err = createUser(jane, uow)
		Ω(err).To(Succeed())

		cmdA = &command.CreateTeam{
			Name:        "Team A",
			Description: "Team A Description",
			User:        john,
		}
		createTeam(cmdA, john)

		cmdB = &command.CreateTeam{
			Name:        "Team B",
			Description: "Team B Description",
			User:        john,
		}
		createTeam(cmdB, john)
	})
	Context("Find a team by ID", func() {
		It("Found", func() {
			team, err := view.Team(cmdA.TeamID, john, Bus.UoW)
			Ω(err).To(Succeed())
			Ω(team.IsPersonal).To(BeFalse())
			Ω(team.Memberships).To(HaveLen(1))
		})
		It("Not Found", func() {
			_, err := view.Team(uuid.NewV4(), john, Bus.UoW)
			Ω(err).To(HaveOccurred())
		})
	})
	It("Update", func() {
		ctx := context.Background()
		cmdUpdate := command.UpdateTeam{
			TeamID:      cmdA.TeamID,
			Name:        "Team C",
			Description: "Team C Description",
			User:        john,
		}
		err := Bus.Handle(ctx, &cmdUpdate)
		Ω(err).To(Succeed())

		team, err := view.Team(cmdA.TeamID, john, Bus.UoW)
		Ω(err).To(Succeed())
		Ω(team.Name).To(Equal("Team C"))
		Ω(team.Description).To(Equal("Team C Description"))
	})
	Context("Get list of teams", func() {
		It("List", func() {
			respPaginated, err := view.Teams(Bus.UoW, john, "", 1, 10)
			Ω(err).To(Succeed())
			Ω(respPaginated.Data).To(HaveLen(3))
			Ω(respPaginated.Page).To(Equal(1))
			Ω(respPaginated.PageSize).To(Equal(10))
			Ω(respPaginated.TotalPage).To(Equal(1))
			Ω(int(respPaginated.TotalData)).To(Equal(3))
			Ω(respPaginated.HasNext).To(BeFalse())
			Ω(respPaginated.HasPrev).To(BeFalse())
		})
	})
	Context("Invite member to a team", Ordered, func() {
		var janeTeam *command.CreateTeam
		BeforeEach(func() {
			janeTeam = &command.CreateTeam{
				TeamID:      uuid.NewV4(),
				Name:        "Team Jane",
				Description: "Team Jane Description",
				User:        jane,
			}
			createTeam(janeTeam, jane)
		})
		It("Invite member", func() {
			ctx := context.Background()
			cmd := command.InviteMember{
				TeamID: janeTeam.TeamID,
				Invitees: []command.Invitee{
					{
						Email: "james@mail.com",
						Role:  domain.Member,
					},
				},
				Sender: jane,
			}
			err := Bus.Handle(ctx, &cmd)
			Ω(err).To(Succeed())

			var invitation domain.Invitation
			row := Bus.UoW.GetDB().QueryRow(
				context.Background(),
				`SELECT i.id, i.email, i.expires_at, i.status, i.team_id, i.role_id, 
				i.sender_id, i.is_active, i.created_at, i.updated_at FROM invitations i 
				WHERE team_id = $1 AND email = $2`, janeTeam.TeamID, "james@mail.com")
			err = row.Scan(
				&invitation.ID,
				&invitation.Email,
				&invitation.ExpiresAt,
				&invitation.Status,
				&invitation.TeamID,
				&invitation.RoleID,
				&invitation.SenderID,
				&invitation.IsActive,
				&invitation.CreatedAt,
				&invitation.UpdatedAt,
			)
			Ω(err).To(Succeed())
			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
		})
		It("Verify invitation", func() {
			ctx := context.Background()
			cmd := command.InviteMember{
				TeamID: janeTeam.TeamID,
				Invitees: []command.Invitee{
					{
						Email: "james@mail.com",
						Role:  domain.Member,
					},
				},
				Sender: jane,
			}
			err := Bus.Handle(ctx, &cmd)
			Ω(err).To(Succeed())

			var invitation domain.Invitation
			row := Bus.UoW.GetDB().QueryRow(
				ctx,
				`SELECT i.id, i.email, i.expires_at, i.status, i.team_id, i.role_id, 
				i.sender_id, i.is_active, i.created_at, i.updated_at FROM invitations i 
				WHERE team_id = $1 AND email = $2`, janeTeam.TeamID, "james@mail.com")
			err = row.Scan(
				&invitation.ID,
				&invitation.Email,
				&invitation.ExpiresAt,
				&invitation.Status,
				&invitation.TeamID,
				&invitation.RoleID,
				&invitation.SenderID,
				&invitation.IsActive,
				&invitation.CreatedAt,
				&invitation.UpdatedAt,
			)
			Ω(err).To(Succeed())

			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
			Ω(invitation.Status).To(Equal(domain.InvitationStatusPending))

			invitation.Status = domain.InvitationStatusSent
			tx, err := Bus.UoW.Begin(ctx)
			Ω(err).To(Succeed())

			err = Bus.UoW.Invitation.Update(invitation, tx)
			Ω(err).To(Succeed())
			tx.Commit(ctx)

			row = Bus.UoW.GetDB().QueryRow(
				ctx,
				`SELECT i.id, i.email, i.expires_at, i.status, i.team_id, i.role_id, 
				i.sender_id, i.is_active, i.created_at, i.updated_at FROM invitations i 
				WHERE team_id = $1 AND email = $2`, janeTeam.TeamID, "james@mail.com")
			err = row.Scan(
				&invitation.ID,
				&invitation.Email,
				&invitation.ExpiresAt,
				&invitation.Status,
				&invitation.TeamID,
				&invitation.RoleID,
				&invitation.SenderID,
				&invitation.IsActive,
				&invitation.CreatedAt,
				&invitation.UpdatedAt,
			)
			Ω(err).To(Succeed())
			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
			Ω(invitation.Status).To(Equal(domain.InvitationStatusSent))

			james := domain.NewUser("James", "Doe", "james@mail.com", "", "Google", true)
			err = createUser(james, Bus.UoW)
			Ω(err).To(Succeed())

			cmdVerify := command.UpdateInvitationStatus{
				InvitationID: invitation.ID,
				Status:       "accepted",
				User:         james,
			}
			err = Bus.Handle(ctx, &cmdVerify)
			Ω(err).To(Succeed())

			row = Bus.UoW.GetDB().QueryRow(
				ctx,
				`SELECT i.id, i.email, i.expires_at, i.status, i.team_id, i.role_id,
				i.sender_id, i.is_active, i.created_at, i.updated_at FROM invitations i
				WHERE team_id = $1 AND email = $2`, janeTeam.TeamID, "james@mail.com")
			err = row.Scan(
				&invitation.ID,
				&invitation.Email,
				&invitation.ExpiresAt,
				&invitation.Status,
				&invitation.TeamID,
				&invitation.RoleID,
				&invitation.SenderID,
				&invitation.IsActive,
				&invitation.CreatedAt,
				&invitation.UpdatedAt,
			)
			Ω(err).To(Succeed())
			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
			Ω(invitation.Status).To(Equal(domain.InvitationStatusAccepted))
		})
	})
})
