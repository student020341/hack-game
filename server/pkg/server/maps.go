package server

import (
	"fmt"
	"net/http"
	"server/pkg/models"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

func (s *Server) listTowns(c echo.Context) error {
	var towns []models.Town
	tx := s.DB.Find(&towns)
	if tx.Error != nil {
		return c.String(http.StatusInternalServerError, "failed to retrieve server towns: "+tx.Error.Error())
	}

	return c.JSON(http.StatusOK, towns)
}

func (s *Server) joinTown(c echo.Context) error {
	account := c.Get("account").(*models.Account)

	// ensure account not already present in a town
	{
		existingPlayer, existingTown, _ := s.findPlayerInServers(account.ID)
		if existingPlayer != nil {
			return c.String(
				http.StatusBadRequest,
				fmt.Sprintf(
					"player %s is already logged into server town '%s' with character %s",
					account.ID,
					existingTown.Name,
					existingPlayer.Character.Name,
				),
			)
		}
	}

	// validate input
	var input struct {
		TownID      *string `json:"townID"`
		CharacterID *string `json:"characterID"`
	}
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "invalid input")
	}

	if input.TownID == nil || input.CharacterID == nil {
		return c.String(http.StatusBadRequest, "townID and characterID required")
	}

	// get character
	var char models.Character
	tx := s.DB.Take(&char, "id = ?", *input.CharacterID)
	if tx.Error != nil {
		return c.String(http.StatusInternalServerError, "failed to retrieve character: "+tx.Error.Error())
	}

	// get town & add user to server collection of players
	for _, town := range s.Towns {
		if town.ID == *input.TownID {
			// TODO create socket or something, checking if player is in server town players
			// ensure player not already in town
			for _, player := range town.Players {
				if player.Account.ID == account.ID {
					return c.String(http.StatusBadRequest, "player has already joined this world")
				}
			}
			town.Players = append(town.Players, models.Player{
				Account:   *account,
				Character: char,
			})
			return c.NoContent(http.StatusOK)
		}
	}

	return c.String(http.StatusInternalServerError, "failed to join town")
}

// reuse?
var upgrader = websocket.Upgrader{}

func (s *Server) joinTownSocket(c echo.Context) error {
	// find which town the player is in
	account := c.Get("account").(*models.Account)
	existingPlayer, existingTown, _ := s.findPlayerInServers(account.ID)
	if existingPlayer == nil {
		return c.String(http.StatusBadRequest, "no player data to join, join a server town first")
	}

	// establish connection
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	// defer ws.Close()

	// store socket in player instance
	for tIdx, t := range s.Towns {
		if t.ID == existingTown.ID {
			for pIdx, p := range t.Players {
				if p.Account.ID == account.ID {
					s.Towns[tIdx].Players[pIdx].Socket = ws
				}
			}
		}
	}

	// greet player
	err = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		return err
	}

	// for {
	// 	// write
	// 	err := ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	// 	if err != nil {
	// 		if err == websocket.ErrCloseSent {
	// 			break
	// 		}
	// 		c.Logger().Error(err)
	// 	}

	// 	// read
	// 	_, msg, err := ws.ReadMessage()
	// 	if err != nil {
	// 		if err == websocket.ErrCloseSent || err.Error() == "websocket: close 1005 (no status)" {
	// 			break
	// 		}
	// 		c.Logger().Error(err)
	// 	}
	// 	fmt.Printf("message: %s\n", msg)
	// }

	return nil
}

func (s *Server) leaveTown(c echo.Context) error {
	account := c.Get("account").(*models.Account)

	// find account
	_, t, index := s.findPlayerInServers(account.ID)
	if t != nil {
		// remove account from town
		t.Players = append(t.Players[:*index], t.Players[*index+1:]...)
		return c.NoContent(http.StatusOK)
	}

	return c.String(http.StatusBadRequest, "account not found")
}

func (s *Server) listTownPlayers(c echo.Context) error {
	townID := c.Param("id")
	var town *models.Town
	for _, t := range s.Towns {
		if t.ID == townID {
			town = t
			break
		}
	}

	if town == nil {
		return c.String(http.StatusBadRequest, "unknown town "+townID)
	}

	return c.JSON(http.StatusOK, town.Players)
}

func (s *Server) findPlayerInServers(id string) (*models.Player, *models.Town, *int) {
	for _, t := range s.Towns {
		for i, p := range t.Players {
			if p.Account.ID == id {
				return &p, t, &i
			}
		}
	}

	return nil, nil, nil
}
