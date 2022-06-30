package integration_test

import (
	"net/http/httptest"
	"testing"

	dbPkg "server/db"
	"server/pkg/server"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var (
	s          server.Server
	httpServer *httptest.Server
)

var _ = BeforeEach(func() {
	var err error
	db := dbPkg.NewDB("file::memory:?cache=shared")
	s = server.MakeServer(db)
	router := s.MakeRoutes()
	httpServer = httptest.NewServer(router)
	Expect(err).To(Succeed())
})

var _ = AfterEach(func() {
	if httpServer != nil {
		httpServer.Close()
	}
})
