package db_test

import (
	"server/pkg/accounts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("db test", func() {
	It("create account and verify password", func() {

		// create account
		acc, err := accounts.CreateAccount("test", "cake")
		Expect(err).To(Succeed())

		DB.Create(&acc)

		// grab account
		var acc2 accounts.Account
		DB.First(&acc2, "ID = ?", acc.ID)

		// test login
		auth, err := accounts.Login("cheese", acc2)
		Expect(err).To(Succeed())
		Expect(auth).To(BeNil())

		auth, err = accounts.Login("cake", acc2)
		Expect(err).To(Succeed())
		Expect(auth).NotTo(BeNil())

		DB.Create(auth)

		// find account from session
		var acc3 accounts.Account
		DB.First(&acc3, "ID = ?", auth.AccountID)

		Expect(acc3.ID).To(Equal(acc.ID))
	})
})
