package integration

import (
	"authorization/domain"
	"authorization/domain/command"
	"authorization/domain/dto"
	"authorization/infrastructure/persistence"
	"authorization/repository"
	"authorization/service/handlers"
	"authorization/view"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func createUser(ctx context.Context, user domain.User) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	Ω(txErr).To(Succeed())

	defer func() {
		tx.Rollback(ctx)
	}()

	_, err := repository.User.Add(ctx, user, tx)
	if err != nil {
		return err
	}

	ownerRole, err := repository.Role.GetByName(ctx, domain.Owner)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Fatal().Err(err).Msgf("Role with name %s is not exist! Detail: %s", domain.Owner, err.Error())
			return err
		}
		return err
	}

	team := domain.NewTeam(user, ownerRole.ID, "", "", true)

	_, err = repository.Team.Add(ctx, team, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

var _ = Describe("User Testing", func() {
	var (
		johnUserId uuid.UUID
	)
	ctx := context.Background()
	BeforeEach(func() {
		john := domain.NewUser("John", "Doe", "johndoe@example.com", "", "Google", true)
		err := createUser(ctx, john)
		Ω(err).To(Succeed())
		johnUserId = john.ID

	})
	Context("Load", func() {
		It("Found", func() {
			user, err := view.User(ctx, johnUserId)
			Ω(err).To(Succeed())
			Ω(user.FirstName).To(Equal("John"))
			Ω(user.LastName).To(Equal("Doe"))
			Ω(user.Email).To(Equal("johndoe@example.com"))
			Ω(user).To(BeAssignableToTypeOf(&dto.PublicUser{}))
		})
		It("Not Found", func() {
			_, err := view.User(ctx, uuid.NewV4())
			Ω(err).To(HaveOccurred())
		})
	})
	It("List", func() {
		respPaginated, err := view.Users(ctx, 1, 10)
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

			jane := domain.NewUser("Jane", "Doe", "janedoe@example.com", "", "Google", true)
			err := createUser(ctx, jane)
			Ω(err).To(Succeed())

			respPaginated, err := view.Users(ctx, 1, 10)
			Ω(err).To(Succeed())
			Ω(respPaginated.Data).To(HaveLen(2))
		})
		It("Update", func() {
			user, err := repository.User.Get(ctx, johnUserId)
			Ω(err).To(Succeed())
			Ω(user.FirstName).To(Equal("John"))
			Ω(user.PhoneNumber).To(Equal(""))
			Ω(user.IsActive).To(BeTrue())

			cmd := command.UpdateUser{
				FirstName:   "Johnny",
				PhoneNumber: "08123456789",
				User:        user,
			}
			err = handlers.UpdateUser(ctx, &cmd)
			Ω(err).To(Succeed())

			user, _ = repository.User.Get(ctx, johnUserId)
			Ω(user.FirstName).To(Equal("Johnny"))
			Ω(user.PhoneNumber).To(Equal("08123456789"))
		})
	})
	It("Delete", func() {
		user, err := repository.User.Get(ctx, johnUserId)
		Ω(err).To(Succeed())
		cmd := command.DeleteUser{
			User: user,
		}
		err = handlers.DeleteUser(ctx, &cmd)
		Ω(err).To(Succeed())
		_, err = view.User(ctx, johnUserId)
		Ω(err).To(HaveOccurred())
	})
})
