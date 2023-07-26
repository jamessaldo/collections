package integration

import (
	"authorization/domain"
	"authorization/domain/command"
	"authorization/util"
	"authorization/view"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	uuid "github.com/satori/go.uuid"
)

func createTeam(cmd *command.CreateTeam, user domain.User) {
	err := Bus.Handle(cmd)
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
		johnUserId uuid.UUID
		janeUserId uuid.UUID
		john       domain.User
		jane       domain.User
		cmdA       *command.CreateTeam
		cmdB       *command.CreateTeam
	)

	BeforeEach(func() {
		uow := Bus.UoW

		now := util.GetTimestampUTC()
		johnUserId = uuid.NewV4()
		janeUserId = uuid.NewV4()
		john = domain.User{
			ID:        johnUserId,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@example.com",
			Username:  "johndoe",
			Provider:  "Google",
			Verified:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

		err := createUser(john, uow)
		Ω(err).To(Succeed())

		jane = domain.User{
			ID:        janeUserId,
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "janedoe@example.com",
			Username:  "janedoe",
			Provider:  "Google",
			Verified:  true,
			CreatedAt: now,
			UpdatedAt: now,
		}

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
		cmdUpdate := command.UpdateTeam{
			TeamID:      cmdA.TeamID,
			Name:        "Team C",
			Description: "Team C Description",
			User:        john,
		}
		err := Bus.Handle(&cmdUpdate)
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
			err := Bus.Handle(&cmd)
			Ω(err).To(Succeed())

			var invitation domain.Invitation
			err = Bus.UoW.GetDB().QueryRow(context.Background(), "SELECT * FROM invitations WHERE team_id = ? AND email = ?", janeTeam.TeamID, "james@mail.com").Scan(&invitation)
			Ω(err).To(Succeed())
			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
		})
		It("Verify invitation", func() {
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
			err := Bus.Handle(&cmd)
			Ω(err).To(Succeed())

			var invitation *domain.Invitation
			err = Bus.UoW.GetDB().QueryRow(context.Background(), "SELECT * FROM invitations WHERE team_id = ? AND email = ?", janeTeam.TeamID, "james@mail.com").Scan(&invitation)
			Ω(err).To(Succeed())

			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
			Ω(invitation.Status).To(Equal(domain.InvitationStatusPending))

			invitation.Status = domain.InvitationStatusSent
			tx, err := Bus.UoW.Begin(context.Background())
			Ω(err).To(Succeed())
			invitation, _ = Bus.UoW.Invitation.Update(invitation, tx)
			Ω(invitation.Status).To(Equal(domain.InvitationStatusSent))

			now := util.GetTimestampUTC()

			james := domain.User{
				ID:        uuid.NewV4(),
				FirstName: "James",
				LastName:  "Doe",
				Email:     "james@mail.com",
				Username:  "jamesdoe",
				Provider:  "Google",
				Verified:  true,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = createUser(james, Bus.UoW)
			Ω(err).To(Succeed())

			cmdVerify := command.UpdateInvitationStatus{
				InvitationID: invitation.ID,
				Status:       "accepted",
				User:         james,
			}
			err = Bus.Handle(&cmdVerify)
			Ω(err).To(Succeed())

			err = Bus.UoW.GetDB().QueryRow(context.Background(), "SELECT * FROM invitations WHERE team_id = ? AND email = ?", janeTeam.TeamID, "james@mail.com").Scan(&invitation)
			Ω(err).To(Succeed())
			Ω(invitation.TeamID).To(Equal(janeTeam.TeamID))
			Ω(invitation.Email).To(Equal("james@mail.com"))
			Ω(invitation.Status).To(Equal(domain.InvitationStatusAccepted))
		})
	})
})
