package integration

import (
	"authorization/domain/command"
	"authorization/domain/dto"
	"authorization/domain/model"
	"authorization/service"
	"authorization/view"
	"errors"
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func createTeam(user *model.User, uow *service.UnitOfWork, tx *gorm.DB) error {
	_, err := uow.User.Add(user, tx)
	if err != nil {
		return err
	}

	ownerRole, err := uow.Role.Get(model.Owner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatalf("Role with name %s is not exist! Detail: %s", model.Owner, err.Error())
			return err
		}
		return err
	}

	membership := user.AddPersonalTeam(ownerRole)

	_, err = uow.Membership.Add(membership, tx)
	if err != nil {
		return err
	}

	return nil
}

var _ = Describe("Team Testing", func() {
	var (
		johnUserId uuid.UUID
		janeUserId uuid.UUID
		john       *model.User
		jane       *model.User
	)
	BeforeEach(func() {
		uow := Bus.UoW
		tx, txErr := uow.Begin(&gorm.Session{})
		Ω(txErr).To(Succeed())

		defer func() {
			tx.Rollback()
		}()

		now := time.Now()
		johnUserId = uuid.NewV4()
		janeUserId = uuid.NewV4()
		john = &model.User{
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

		err := createUser(john, uow, tx)
		Ω(err).To(Succeed())

		jane = &model.User{
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

		err = createUser(jane, uow, tx)
		Ω(err).To(Succeed())

		tx.Commit()

	})
	Context("Load", func() {
		It("Found", func() {
			teams, err := view.Teams(Bus.UoW, john, "", 1, 10)
			Ω(err).To(Succeed())
			Ω(teams.Data).To(HaveLen(1))
			Ω(teams.Data[0].(*dto.TeamRetrievalSchema).IsPersonal).To(BeTrue())
		})
		It("Not Found", func() {
			_, err := view.Team(uuid.NewV4(), john, Bus.UoW)
			Ω(err).To(HaveOccurred())
		})
	})
	Context("Save", func() {
		It("Create", func() {
			cmd := command.CreateTeam{
				TeamID:      uuid.NewV4(),
				Name:        "Team A",
				Description: "Team A Description",
				User:        john,
			}
			err := Bus.Handle(&cmd)
			Ω(err).To(Succeed())

			team, err := view.Team(cmd.TeamID, john, Bus.UoW)
			Ω(err).To(Succeed())
			Ω(team.Name).To(Equal("Team A"))
			Ω(team.IsPersonal).To(BeFalse())
		})
		It("Update", func() {
			cmd := command.CreateTeam{
				TeamID:      uuid.NewV4(),
				Name:        "Team A",
				Description: "Team A Description",
				User:        john,
			}
			err := Bus.Handle(&cmd)
			Ω(err).To(Succeed())

			cmdUpdate := command.UpdateTeam{
				TeamID:      cmd.TeamID,
				Name:        "Team B",
				Description: "Team B Description",
				User:        john,
			}
			err = Bus.Handle(&cmdUpdate)
			Ω(err).To(Succeed())

			team, err := view.Team(cmd.TeamID, john, Bus.UoW)
			Ω(err).To(Succeed())
			Ω(team.Name).To(Equal("Team B"))
			Ω(team.Description).To(Equal("Team B Description"))
		})
		It("List", func() {
			respPaginated, err := view.Teams(Bus.UoW, john, "", 1, 10)
			Ω(err).To(Succeed())
			fmt.Println("aselole", respPaginated.Data)
			Ω(respPaginated.Data).To(HaveLen(1))
			Ω(respPaginated.Page).To(Equal(1))
			Ω(respPaginated.PageSize).To(Equal(10))
			Ω(respPaginated.TotalPage).To(Equal(1))
			Ω(int(respPaginated.TotalData)).To(Equal(1))
			Ω(respPaginated.HasNext).To(BeFalse())
			Ω(respPaginated.HasPrev).To(BeFalse())
		})
	})
	// It("Delete", func() {
	// 	user, err := Bus.UoW.User.Get(johnUserId)
	// 	Ω(err).To(Succeed())
	// 	cmd := command.DeleteUser{
	// 		User: user,
	// 	}
	// 	err = Bus.Handle(&cmd)
	// 	Ω(err).To(Succeed())
	// 	_, err = view.User(johnUserId, Bus.UoW)
	// 	Ω(err).To(HaveOccurred())
	// })
})
