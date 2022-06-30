package db_test

import (
	"testing"

	dbPkg "server/db"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Suite")
}

var (
	DB *gorm.DB
)

var _ = BeforeEach(func() {
	DB = dbPkg.NewDB("file::memory:?cache=shared")
})
