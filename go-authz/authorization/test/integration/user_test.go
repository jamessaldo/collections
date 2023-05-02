package integration

import (
	"auth/domain/model"
	"auth/infrastructure/worker"
	"auth/service"
	"auth/service/handlers"
	"errors"
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var _ = Describe("Repository", func() {
	var (
		bus *service.MessageBus
	)
	BeforeEach(func() {
		err := Db.AutoMigrate(
			&model.User{},
			&model.Endpoint{},
			&model.Role{},
			&model.Access{},
			&model.Membership{},
			&model.Team{},
			&model.Invitation{})
		if err != nil {
			log.Fatal(err)
		}

		uow, err := service.NewUnitOfWork(Db)
		if err != nil {
			log.Fatal(err)
		}

		asynqClient := worker.CreateAsynqClient()
		defer asynqClient.Close()

		mailer := worker.NewMailer(asynqClient)
		if err != nil {
			log.Fatal(err)
		}

		Ω(err).To(Succeed())

		bus = service.NewMessageBus(handlers.COMMAND_HANDLERS, uow, mailer)

		tx, txErr := uow.Begin(&gorm.Session{})
		if txErr != nil {
			log.Fatal(txErr)
		}

		defer func() {
			tx.Rollback()
		}()

		now := time.Now()

		uow = bus.UoW

		user := &model.User{
			ID:          uuid.NewV4(),
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "johndoe@example.com",
			Password:    "",
			PhoneNumber: "",
			AvatarURL:   "",
			Provider:    "Google",
			Verified:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		_, err = uow.User.Add(user, tx)
		if err != nil {
			log.Fatal(err)
		}

		ownerRole, err := uow.Role.Get(model.Owner)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("Role with name %s is not exist! Detail: %s", model.Owner, err.Error())
			}
			log.Fatal(err)
		}

		membership := user.AddPersonalTeam(ownerRole)

		_, err = uow.Membership.Add(membership, tx)
		if err != nil {
			log.Fatal(err)
		}
		Ω(err).To(Succeed())
	})
	// Context("Load", func() {
	// 	It("Found", func() {
	// 		blog, err := repo.Load(1)

	// 		Ω(err).To(Succeed())
	// 		Ω(blog.Content).To(Equal("hello"))
	// 		Ω(blog.Tags).To(Equal(pq.StringArray{"a", "b"}))
	// 	})
	// 	It("Not Found", func() {
	// 		_, err := repo.Load(999)
	// 		Ω(err).To(HaveOccurred())
	// 	})
	// })
	It("ListAll", func() {
		l, err := bus.UoW.User.List(1, 10)
		Ω(err).To(Succeed())
		Ω(l).To(HaveLen(1))
	})
	// It("List", func() {
	// 	l, err := repo.List(0, 10)
	// 	Ω(err).To(Succeed())
	// 	Ω(l).To(HaveLen(1))
	// })
	// Context("Save", func() {
	// 	It("Create", func() {
	// 		blog := &dbtest.Blog{
	// 			Title:     "post2",
	// 			Content:   "hello",
	// 			Tags:      []string{"a", "b"},
	// 			CreatedAt: time.Now(),
	// 		}
	// 		err := repo.Save(blog)
	// 		Ω(err).To(Succeed())
	// 		Ω(blog.ID).To(BeEquivalentTo(2))
	// 	})
	// 	It("Update", func() {
	// 		blog, err := repo.Load(1)
	// 		Ω(err).To(Succeed())

	// 		blog.Title = "foo"
	// 		err = repo.Save(blog)
	// 		Ω(err).To(Succeed())
	// 	})
	// })
	// It("Delete", func() {
	// 	err := repo.Delete(1)
	// 	Ω(err).To(Succeed())
	// 	_, err = repo.Load(1)
	// 	Ω(err).To(HaveOccurred())
	// })
	// DescribeTable("SearchByTitle",
	// 	func(q string, found int) {
	// 		l, err := repo.SearchByTitle(q, 0, 10)
	// 		Ω(err).To(Succeed())
	// 		Ω(l).To(HaveLen(found))
	// 	},
	// 	Entry("found", "post", 1),
	// 	Entry("not found", "bar", 0),
	// )
})
