package integration

import (
	"authorization/domain"
	"authorization/domain/command"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/worker"
	"authorization/repository"
	"authorization/service/handlers"
	"authorization/view"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	uuid "github.com/satori/go.uuid"
)

func createTeam(ctx context.Context, cmd *command.CreateTeam, user domain.User) {
	err := handlers.CreateTeam(ctx, cmd)
	Ω(err).To(Succeed())

	team, err := view.Team(ctx, cmd.TeamID, user)
	Ω(err).To(Succeed())
	Ω(team.Name).To(Equal(cmd.Name))
	Ω(team.Description).To(Equal(cmd.Description))
	Ω(team.IsPersonal).To(BeFalse())
}

var _ = Describe("Team Testing", Ordered, func() {
	format.MaxLength = 0
	ctx := context.Background()
	client := worker.CreateMailerClientMock()
	worker.CreateMailerMock(client)

	var (
		john domain.User
		jane domain.User
		cmdA *command.CreateTeam
		cmdB *command.CreateTeam
	)

	BeforeEach(func() {
		john = domain.NewUser("John", "Doe", "johndoe@example.com", "", "Google", true)
		err := createUser(ctx, john)
		Ω(err).To(Succeed())

		jane = domain.NewUser("Jane", "Doe", "janedoe@example.com", "", "Google", true)
		err = createUser(ctx, jane)
		Ω(err).To(Succeed())

		cmdA = &command.CreateTeam{
			Name:        "Team A",
			Description: "Team A Description",
			User:        john,
		}
		createTeam(ctx, cmdA, john)

		cmdB = &command.CreateTeam{
			Name:        "Team B",
			Description: "Team B Description",
			User:        john,
		}
		createTeam(ctx, cmdB, john)
	})
	Context("Find a team by ID", func() {
		It("Found", func() {
			team, err := view.Team(ctx, cmdA.TeamID, john)
			Ω(err).To(Succeed())
			Ω(team.IsPersonal).To(BeFalse())
			Ω(team.Memberships).To(HaveLen(1))
		})
		It("Not Found", func() {
			_, err := view.Team(ctx, uuid.NewV4(), john)
			Ω(err).To(HaveOccurred())
		})
	})
	It("Update", func() {
		cmdUpdate := command.UpdateTeam{
			TeamID:      cmdA.TeamID,
			Name:        "Team C",
			Description: "Team C Description",
			User:        john,
		}
		err := handlers.UpdateTeam(ctx, &cmdUpdate)
		Ω(err).To(Succeed())

		team, err := view.Team(ctx, cmdA.TeamID, john)
		Ω(err).To(Succeed())
		Ω(team.Name).To(Equal("Team C"))
		Ω(team.Description).To(Equal("Team C Description"))
	})
	Context("Get list of teams", func() {
		It("List", func() {
			respPaginated, err := view.Teams(ctx, john, "", 1, 10)
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
			createTeam(ctx, janeTeam, jane)
		})
		It("Invite member", func() {
			cmd := command.SendInvitation{
				TeamID: janeTeam.TeamID,
				Invitees: []command.Invitee{
					{
						Email: "james@mail.com",
						Role:  domain.Member,
					},
				},
				Sender: jane,
			}
			err := handlers.SendInvitation(ctx, &cmd)
			Ω(err).To(Succeed())

			var invitation domain.Invitation
			row := persistence.Pool.QueryRow(
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
			cmd := command.SendInvitation{
				TeamID: janeTeam.TeamID,
				Invitees: []command.Invitee{
					{
						Email: "james@mail.com",
						Role:  domain.Member,
					},
				},
				Sender: jane,
			}
			err := handlers.SendInvitation(ctx, &cmd)
			Ω(err).To(Succeed())

			var invitation domain.Invitation
			row := persistence.Pool.QueryRow(
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
			tx, err := persistence.Pool.Begin(ctx)
			Ω(err).To(Succeed())

			err = repository.Invitation.Update(ctx, invitation, tx)
			Ω(err).To(Succeed())
			tx.Commit(ctx)

			row = persistence.Pool.QueryRow(
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
			err = createUser(ctx, james)
			Ω(err).To(Succeed())

			cmdVerify := command.UpdateInvitationStatus{
				InvitationID: invitation.ID,
				Status:       "accepted",
				User:         james,
			}

			err = handlers.UpdateInvitationStatus(ctx, &cmdVerify)
			Ω(err).To(Succeed())

			row = persistence.Pool.QueryRow(
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
