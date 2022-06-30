package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"server/db"
	"server/pkg/accounts"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

// enforce login is admin
func AdminAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		vagueNegativeResposne := func() error {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "auth :)"})
		}

		if token == "" {
			return vagueNegativeResposne()
		}

		// check db
		var authThing *accounts.AuthSession
		db.DB.Find(&authThing, "Token = ?", token)
		if authThing == nil {
			return vagueNegativeResposne()
		}

		var acc accounts.Account
		db.DB.Find(&acc, "ID = ?", authThing.AccountID)
		if acc.ID == "" || acc.Level != 0 {
			return vagueNegativeResposne()
		}

		return next(c)
	}
}

// keep track of when a login token was last used
func UpdateAccessTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("token")
		if token != "" {
			var authThing accounts.AuthSession
			tx := db.DB.Find(&authThing, "Token = ?", token)
			if tx.Error == nil && authThing.AccountID != "" {
				authThing.LastAccessed = time.Now()
				db.DB.Save(&authThing)
			}
		}

		return next(c)
	}
}

func main() {
	db.Init()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(UpdateAccessTime)

	// admin stuff
	e.GET("/api/accounts", listAccounts, AdminAuth)

	// open
	e.POST("/api/accounts", createAccount)
	e.POST("/api/login", login)
	e.POST("/api/logout", logout)

	// test
	var stuff []accounts.AuthSession
	db.DB.Find(&stuff)
	jstr, _ := json.MarshalIndent(stuff, "", "  ")
	fmt.Printf("\n\n%s\n\n", jstr)

	e.Logger.Fatal(e.Start(":2001"))
}

func listAccounts(c echo.Context) error {
	var accounts []accounts.Account
	db.DB.Find(&accounts)

	return c.JSON(http.StatusOK, accounts)
}

func createAccount(c echo.Context) error {
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

	tx := db.DB.Create(&acc)
	if tx.Error != nil {
		log.Errorf("error saving account: %+v", tx.Error)
		return c.String(http.StatusInternalServerError, tx.Error.Error())
	}

	return c.NoContent(http.StatusOK)
}

func login(c echo.Context) error {
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
	tx := db.DB.First(&acc, "Username = ?", *input.Username)
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

	tx = db.DB.Create(&authThing)
	if tx.Error != nil {
		return c.String(http.StatusInternalServerError, "failed to save login session: "+tx.Error.Error())
	}

	return c.String(http.StatusOK, authThing.Token)
}

func logout(c echo.Context) error {
	token := c.Request().Header.Get("token")
	if token == "" {
		return c.String(http.StatusBadRequest, "no session token")
	}

	all := c.Request().URL.Query().Get("all")

	var auth accounts.AuthSession
	tx := db.DB.Take(&auth, "Token = ?", token)
	if tx.Error != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if all == "1" {
		tx = db.DB.Delete(&accounts.AuthSession{}, "account_id = ?", auth.AccountID)
	} else {
		tx = db.DB.Delete(&auth)
	}

	if tx.Error != nil {
		log.Errorf("failed to delete auth: %+v", tx.Error)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}
