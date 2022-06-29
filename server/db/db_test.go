package db_test

import (
	"server/pkg/accounts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("db test", func() {
	It("create account and verify password", func() {

		acc, err := accounts.CreateAccount("test", "cake")
		Expect(err).To(Succeed())

		DB.Create(&acc)

		var acc2 accounts.Account
		DB.First(&acc2)

		ok, err := accounts.VerifyLogin("cheese", acc2)
		Expect(err).To(Succeed())
		Expect(ok).To(BeFalse())

		ok, err = accounts.VerifyLogin("cake", acc2)
		Expect(err).To(Succeed())
		Expect(ok).To(BeTrue())
	})
})
