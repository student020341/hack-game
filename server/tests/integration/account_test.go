package integration_test

import (
	"fmt"

	"server/pkg/accounts"
	testPkg "server/tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("account tests", func() {
	It("sanity test", func() {
		code, body, err := testPkg.SimpleGet(fmt.Sprintf("%s/status", httpServer.URL))
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
		code, _, err := testPkg.SimplePost(fmt.Sprintf("%s/api/accounts", httpServer.URL), input)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))

		// verify account exists
		var acc accounts.Account
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
		code, token, err := testPkg.SimplePost(url, input)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))

		// verify log in token created
		var auth accounts.AuthSession
		tx := s.DB.Take(&auth, "Token = ?", string(token))
		Expect(tx.Error).To(Succeed())
		Expect(auth.Token).To(Equal(string(token)))
	})
})
