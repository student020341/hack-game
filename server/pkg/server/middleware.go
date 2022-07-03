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
		account, _ := c.Get("account").(*models.Account)

		if account == nil || account.Level != 0 {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

// ensure some user is logged in
func (s *Server) AnyAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		account, _ := c.Get("account").(*models.Account)

		if account == nil || account.ID == "" {
			return c.NoContent(http.StatusUnauthorized)
		}

		return next(c)
	}
}

// looks up account if token header is present and stores account in context
func (s *Server) SetAccountFromToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		if token == "" {
			return next(c)
		}

		var auth models.AuthSession
		s.DB.Take(&auth, "Token = ?", token)
		if auth.AccountID == "" {
			return next(c)
		}

		var acc models.Account
		s.DB.Take(&acc, "ID = ?", auth.AccountID)
		if acc.ID != "" {
			c.Set("account", &acc)
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
