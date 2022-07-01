package integration_test

import (
	"fmt"

	accountPkg "server/pkg/accounts"
	testPkg "server/tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("account tests", func() {
	It("sanity test", func() {
		code, body, err := testPkg.SimpleGet(fmt.Sprintf("%s/status", httpServer.URL), nil)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))
		Expect(string(body)).To(Equal("ok"))
	})

	It("should create an account", func() {
		// create account
		input := struct {
			Username string
			Password string
		}{
			Username: "test",
			Password: "foo",
		}
		code, _, err := testPkg.SimplePost(fmt.Sprintf("%s/api/accounts", httpServer.URL), input, nil)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))

		// verify account exists
		var acc accountPkg.Account
		tx := s.DB.Take(&acc, "Username = ?", input.Username)
		Expect(tx.Error).To(Succeed())
		Expect(acc.Username).To(Equal(input.Username))
	})

	It("should log in to existing account", func() {
		input := struct {
			Username string
			Password string
		}{
			Username: "admin",
			Password: "admin",
		}

		url := fmt.Sprintf("%s/api/login", httpServer.URL)
		code, token, err := testPkg.SimplePost(url, input, nil)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))

		// verify log in token created
		var auth accountPkg.AuthSession
		tx := s.DB.Take(&auth, "Token = ?", string(token))
		Expect(tx.Error).To(Succeed())
		Expect(auth.Token).To(Equal(string(token)))
	})

	FContext("logged in accounts", func() {
		var acc *accountPkg.Account
		var auth *accountPkg.AuthSession
		BeforeEach(func() {
			// create user and session
			var err error
			acc, err = accountPkg.CreateAccount("user", "pass")
			Expect(err).To(Succeed())

			auth, err = accountPkg.Login("pass", *acc)
			Expect(err).To(Succeed())

			// persist
			tx := s.DB.Create(&acc)
			Expect(tx.Error).To(Succeed())
			tx = s.DB.Create(&auth)
			Expect(tx.Error).To(Succeed())
		})

		It("user can log out", func() {
			code, _, err := testPkg.SimplePost(httpServer.URL+"/api/logout", nil, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify session deleted
			tx := s.DB.Take(&accountPkg.AuthSession{}, "Token = ?", &auth.Token)
			Expect(tx.Error).To(Equal(gorm.ErrRecordNotFound))
		})

		It("user can log out of all sessions", func() {
			// create a few more
			for i := 0; i < 5; i++ {
				a, err := accountPkg.Login("pass", *acc)
				Expect(err).To(Succeed())
				tx := s.DB.Create(&a)
				Expect(tx.Error).To(Succeed())
			}

			// verify there are several sessions
			var count int64
			tx := s.DB.Find(&accountPkg.AuthSession{}).Count(&count)
			Expect(tx.Error).To(Succeed())
			Expect(int(count)).To(Equal(6))

			code, _, err := testPkg.SimplePost(httpServer.URL+"/api/logout?all=1", nil, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify there are 0 sessions
			tx = s.DB.Find(&accountPkg.AuthSession{}).Count(&count)
			Expect(tx.Error).To(Succeed())
			Expect(int(count)).To(Equal(0))
		})

		It("regular user cannot access account list", func() {
			code, _, err := testPkg.SimpleGet(fmt.Sprintf("%s/api/accounts", httpServer.URL), nil)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(401))
		})
	}) // logged in accounts
})
