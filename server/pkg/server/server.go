package server

import (
	"net/http"
	dbPkg "server/db"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	DB *gorm.DB
}

func (s *Server) Start() {
	router := s.MakeRoutes()

	router.Logger.Fatal(router.Start(":2001"))
}

func (s *Server) MakeRoutes() *echo.Echo {
	e := echo.New()

	// "github.com/labstack/echo/v4/middleware"
	// e.Use(middleware.Logger())
	e.Use(s.SetAccountFromToken)

	// admin stuff
	e.GET("/api/accounts", s.listAccounts, s.AdminAuth)

	// open
	e.GET("/status", s.getStatus)
	e.POST("/api/accounts", s.createAccount)
	e.POST("/api/login", s.login, s.UpdateAccessTime)
	e.POST("/api/logout", s.logout)

	// logged in
	e.GET("/api/characters", s.getCharacter, s.AnyAuth)
	e.POST("/api/characters", s.createCharacter, s.AnyAuth)

	return e
}

func MakeServer(db *gorm.DB) Server {
	var DB *gorm.DB
	if db != nil {
		DB = db
	} else {
		DB = dbPkg.NewDB("test.db")
	}

	return Server{
		DB: DB,
	}
}

// misc routes
func (s *Server) getStatus(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
