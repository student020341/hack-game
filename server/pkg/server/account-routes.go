package server

import (
	"fmt"
	"log"
	"net/http"
	"server/pkg/accounts"
	"server/pkg/models"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

	if len(*input.Username) == 0 {
		return c.String(http.StatusBadRequest, "username cannot be blank")
	}

	if len(*input.Password) == 0 {
		return c.String(http.StatusBadRequest, "password cannot be blank")
	}

	acc, err := accounts.CreateAccount(*input.Username, *input.Password)
	if err != nil {
		log.Printf("error creating account: %+v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	tx := s.DB.Create(&acc)
	if tx.Error != nil {
		log.Printf("error saving account: %+v", tx.Error)
		return c.String(http.StatusInternalServerError, tx.Error.Error())
	}

	auth, err := accounts.Login(*input.Password, *acc)
	if err != nil {
		return c.String(http.StatusInternalServerError, "account created but failed to login")
	}

	return c.String(http.StatusOK, auth.Token)
}

func (s *Server) login(c echo.Context) error {
	account, _ := c.Get("account").(*models.Account)
	// skip login if user is already logged in
	if account != nil {
		return c.NoContent(http.StatusNoContent)
	}

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

	var acc models.Account
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

	var auth models.AuthSession
	tx := s.DB.Take(&auth, "Token = ?", token)
	if tx.Error != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if all == "1" {
		tx = s.DB.Delete(&models.AuthSession{}, "account_id = ?", auth.AccountID)
	} else {
		tx = s.DB.Delete(&auth)
	}

	if tx.Error != nil {
		fmt.Printf("failed to delete auth: %+v", tx.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

// TODO move these to character file?
func (s *Server) getCharacter(c echo.Context) error {
	// route should not be reachable without logging in due to middleware, panic
	account := c.Get("account").(*models.Account)

	var chars []models.Character
	err := s.DB.Model(&account).Association("Characters").Find(&chars)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, chars)
}

func (s *Server) createCharacter(c echo.Context) error {
	account := c.Get("account").(*models.Account)

	var input struct {
		Name *string
	}

	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "invalid input")
	}

	if input.Name == nil {
		return c.String(http.StatusBadRequest, "need character name")
	}

	// add character to account
	character := models.Character{
		ID:   uuid.New().String(),
		Name: *input.Name,
	}
	err := s.DB.Model(&account).Association("Characters").Append(&character)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (s *Server) deleteCharacter(c echo.Context) error {
	account := c.Get("account").(*models.Account)

	var input struct {
		ID *string `json:"id"`
	}

	err := c.Bind(&input)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid input")
	}

	if input.ID == nil {
		return c.String(http.StatusBadRequest, "input 'ID' is required")
	}

	// get
	var char models.Character
	err = s.DB.Model(&account).Association("Characters").Find(&char, "id = ?", *input.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, "error retrieving character: "+err.Error())
	}

	if char.ID == "" {
		return c.String(http.StatusInternalServerError, "unknown error retrieving character")
	}

	// delete
	tx := s.DB.Delete(&char)
	if tx.Error != nil {
		return c.String(http.StatusInternalServerError, "error deleting character: "+tx.Error.Error())
	}

	return c.NoContent(http.StatusOK)
}
