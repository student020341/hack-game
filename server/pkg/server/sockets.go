package server

import (
	"fmt"
	"server/pkg/models"

	"github.com/gorilla/websocket"
)

func (s *Server) handlePlayerSocket(ws *websocket.Conn, townID string, playerID string) {
	var town *models.Town
	for i, t := range s.Towns {
		if t.ID == townID {
			town = s.Towns[i]
		}
	}

	// TODO error handling
	if town == nil {
		return
	}

	for {
		// listen for player messages until socket disconnects
		_, msg, err := ws.ReadMessage()
		if err != nil {
			// TODO handle errors
			// TODO use logger instead
			fmt.Printf("player %s socket error: %s\n", playerID, err.Error())
			break
		}
		town.MessageChannel <- models.PlayerMessage{ID: playerID, Message: string(msg)}
	}

	// clear socket
	for tIndex, t := range s.Towns {
		if t.ID == townID {
			for pIndex, p := range t.Players {
				if p.Account.ID == playerID && s.Towns[tIndex].Players[pIndex].Socket != nil {
					s.Towns[tIndex].Players[pIndex].Socket.Close()
					s.Towns[tIndex].Players[pIndex].Socket = nil
				}
			}
		}
	}
}
