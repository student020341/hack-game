package server

import (
	"log"
	"net/http"
	dbPkg "server/db"
	"server/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	DB    *gorm.DB
	Towns []*models.Town
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
	e.DELETE("/api/characters", s.deleteCharacter, s.AnyAuth)

	// server towns
	e.GET("/api/servers", s.listTowns, s.AnyAuth)
	e.POST("/api/servers/join", s.joinTown, s.AnyAuth)
	e.GET("/api/servers/socket", s.joinTownSocket, s.AnyAuth)
	e.GET("/api/servers/leave", s.leaveTown, s.AnyAuth)
	e.GET("/api/servers/:id/players", s.listTownPlayers, s.AnyAuth)

	return e
}

func MakeServer(db *gorm.DB) Server {
	var DB *gorm.DB
	if db != nil {
		DB = db
	} else {
		DB = dbPkg.NewDB("test.db")
	}

	// load towns
	var towns []*models.Town
	tx := DB.Find(&towns)
	if tx.Error != nil {
		log.Fatalf("failed to load towns: %+v", tx.Error)
	}

	// TODO should this be here?
	// handle sockets
	for _, t := range towns {
		// initialize channels
		t.MessageChannel = make(chan models.PlayerMessage)
		go func(town *models.Town) {
			for {
				msg := <-town.MessageChannel
				for _, p := range town.Players {
					if p.Account.ID != msg.ID {
						p.Socket.WriteMessage(websocket.TextMessage, []byte(msg.Message))
					}
				}
			}
		}(t)
	}

	return Server{
		DB:    DB,
		Towns: towns,
	}
}

// misc routes
func (s *Server) getStatus(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
