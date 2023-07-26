package integration

import (
	"authorization/domain"
	"authorization/domain/command"
	"authorization/domain/dto"
	"authorization/service"
	"authorization/util"
	"authorization/view"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func createUser(user domain.User, uow *service.UnitOfWork) error {
	ctx := context.Background()
	tx, txErr := uow.Begin(ctx)
	Ω(txErr).To(Succeed())

	defer func() {
		tx.Rollback(ctx)
	}()

	_, err := uow.User.Add(user, tx)
	if err != nil {
		return err
	}

	ownerRole, err := uow.Role.Get(ctx, domain.Owner)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Fatal().Err(err).Msgf("Role with name %s is not exist! Detail: %s", domain.Owner, err.Error())
			return err
		}
		return err
	}

	team := domain.NewTeam(user, ownerRole.ID, "", "", true)

	_, err = uow.Team.Add(team, tx)
	if err != nil {
		return err
	}
	tx.Commit(ctx)

	return nil
}

var _ = Describe("User Testing", func() {
	var (
		johnUserId uuid.UUID
	)
	BeforeEach(func() {
		uow := Bus.UoW

		now := util.GetTimestampUTC()
		johnUserId = uuid.NewV4()
		user := domain.User{
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

		err := createUser(user, uow)
		Ω(err).To(Succeed())

	})
	Context("Load", func() {
		It("Found", func() {
			user, err := view.User(johnUserId, Bus.UoW)
			Ω(err).To(Succeed())
			Ω(user.FirstName).To(Equal("John"))
			Ω(user.LastName).To(Equal("Doe"))
			Ω(user.Email).To(Equal("johndoe@example.com"))
			Ω(user).To(BeAssignableToTypeOf(&dto.PublicUser{}))
		})
		It("Not Found", func() {
			_, err := view.User(uuid.NewV4(), Bus.UoW)
			Ω(err).To(HaveOccurred())
		})
	})
	It("List", func() {
		respPaginated, err := view.Users(Bus.UoW, 1, 10)
		Ω(err).To(Succeed())
		Ω(respPaginated.Data).To(HaveLen(1))
		Ω(respPaginated.Page).To(Equal(1))
		Ω(respPaginated.PageSize).To(Equal(10))
		Ω(respPaginated.TotalPage).To(Equal(1))
		Ω(int(respPaginated.TotalData)).To(Equal(1))
		Ω(respPaginated.HasNext).To(BeFalse())
		Ω(respPaginated.HasPrev).To(BeFalse())
	})
	Context("Save", func() {
		It("Create", func() {
			uow := Bus.UoW

			now := util.GetTimestampUTC()
			janeUserId := uuid.NewV4()
			user := domain.User{
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

			err := createUser(user, uow)
			Ω(err).To(Succeed())

			respPaginated, err := view.Users(Bus.UoW, 1, 10)
			Ω(err).To(Succeed())
			Ω(respPaginated.Data).To(HaveLen(2))
		})
		It("Update", func() {
			user, err := Bus.UoW.User.Get(johnUserId)
			Ω(err).To(Succeed())
			Ω(user.FirstName).To(Equal("John"))
			Ω(user.PhoneNumber).To(Equal(""))
			Ω(user.IsActive).To(BeTrue())

			cmd := command.UpdateUser{
				FirstName:   "Johnny",
				PhoneNumber: "08123456789",
				User:        user,
			}
			err = Bus.Handle(&cmd)
			Ω(err).To(Succeed())

			user, _ = Bus.UoW.User.Get(johnUserId)
			Ω(user.FirstName).To(Equal("Johnny"))
			Ω(user.PhoneNumber).To(Equal("08123456789"))
		})
	})
	It("Delete", func() {
		user, err := Bus.UoW.User.Get(johnUserId)
		Ω(err).To(Succeed())
		cmd := command.DeleteUser{
			User: user,
		}
		err = Bus.Handle(&cmd)
		Ω(err).To(Succeed())
		_, err = view.User(johnUserId, Bus.UoW)
		Ω(err).To(HaveOccurred())
	})
})
