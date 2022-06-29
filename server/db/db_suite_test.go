package db_test

import (
	"testing"

	"server/db"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Suite")
}

var DB *gorm.DB

var _ = BeforeEach(func() {
	var err error
	DB = db.NewDB("file::memory:?cache=shared")
	Expect(err).To(Succeed())
})
