package server

import (
	"net/http"
	"server/pkg/models"

	"github.com/labstack/echo/v4"
)

func (s *Server) listAccounts(c echo.Context) error {
	var accounts []models.Account
	s.DB.Find(&accounts)

	return c.JSON(http.StatusOK, accounts)
}
