package server

import (
	"net/http"
	"server/pkg/accounts"

	"github.com/labstack/echo/v4"
	"github.com/nrfta/go-log"
)

func (s *Server) createAccount(c echo.Context) error {
	var input struct {
		Username *string
		Password *string
	}

	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "invalid input")
	}

	if input.Username == nil || input.Password == nil {
		return c.String(http.StatusBadRequest, "need username and password")
	}

	acc, err := accounts.CreateAccount(*input.Username, *input.Password)
	if err != nil {
		log.Errorf("error creating account: %+v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	tx := s.DB.Create(&acc)
	if tx.Error != nil {
		log.Errorf("error saving account: %+v", tx.Error)
		return c.String(http.StatusInternalServerError, tx.Error.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) login(c echo.Context) error {
	var input struct {
		Username *string
		Password *string
	}

	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "invalid input")
	}

	if input.Username == nil || input.Password == nil {
		return c.String(http.StatusBadRequest, "need username and password")
	}

	var acc accounts.Account
	tx := s.DB.First(&acc, "Username = ?", *input.Username)
	if tx.Error != nil {
		return c.String(http.StatusInternalServerError, "failed to look up account: "+tx.Error.Error())
	}

	authThing, err := accounts.Login(*input.Password, acc)
	if err != nil {
		return c.String(http.StatusInternalServerError, "login error: "+err.Error())
	}

	if authThing == nil {
		return c.NoContent(http.StatusNotFound)
	}

	tx = s.DB.Create(&authThing)
	if tx.Error != nil {
		return c.String(http.StatusInternalServerError, "failed to save login session: "+tx.Error.Error())
	}

	return c.String(http.StatusOK, authThing.Token)
}

func (s *Server) logout(c echo.Context) error {
	token := c.Request().Header.Get("token")
	if token == "" {
		return c.String(http.StatusBadRequest, "no session token")
	}

	all := c.Request().URL.Query().Get("all")

	var auth accounts.AuthSession
	tx := s.DB.Take(&auth, "Token = ?", token)
	if tx.Error != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if all == "1" {
		tx = s.DB.Delete(&accounts.AuthSession{}, "account_id = ?", auth.AccountID)
	} else {
		tx = s.DB.Delete(&auth)
	}

	if tx.Error != nil {
		log.Errorf("failed to delete auth: %+v", tx.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
