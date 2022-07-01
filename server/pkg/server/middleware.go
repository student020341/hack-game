package server

import (
	"net/http"
	"server/pkg/models"
	"time"

	"github.com/labstack/echo/v4"
)

// enforce login is admin
func (s *Server) AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		vagueNegativeResposne := func() error {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "auth :)"})
		}

		if token == "" {
			return vagueNegativeResposne()
		}

		// check db
		var authThing *models.AuthSession
		s.DB.Find(&authThing, "Token = ?", token)
		if authThing == nil {
			return vagueNegativeResposne()
		}

		var acc models.Account
		s.DB.Find(&acc, "ID = ?", authThing.AccountID)
		if acc.ID == "" || acc.Level != 0 {
			return vagueNegativeResposne()
		}

		return next(c)
	}
}

// keep track of when a login token was last used
func (s *Server) UpdateAccessTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		if token != "" {
			var authThing models.AuthSession
			tx := s.DB.Find(&authThing, "Token = ?", token)
			if tx.Error == nil && authThing.AccountID != "" {
				authThing.LastAccessed = time.Now()
				s.DB.Save(&authThing)
			}
		}

		return next(c)
	}
}
