package integration

import (
	"authorization/domain/command"
	"authorization/domain/dto"
	"authorization/domain/model"
	"authorization/service"
	"authorization/view"
	"errors"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func createUser(user *model.User, uow *service.UnitOfWork) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	Ω(txErr).To(Succeed())

	defer func() {
		tx.Rollback()
	}()

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

	team := model.NewTeam(user, uuid.NewV4(), ownerRole.ID, "", "", true)

	_, err = uow.Team.Add(team, tx)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

var _ = Describe("User Testing", func() {
	var (
		johnUserId uuid.UUID
	)
	BeforeEach(func() {
		uow := Bus.UoW

		now := time.Now()
		johnUserId = uuid.NewV4()
		user := &model.User{
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

			now := time.Now()
			janeUserId := uuid.NewV4()
			user := &model.User{
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
