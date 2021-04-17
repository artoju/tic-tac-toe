package socket

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// lobbyClientRead reads messages from lobby connections
func (c *Client) lobbyClientRead() {
	defer func() {
		log.WithFields(log.Fields{
			"playerID": c.ID,
		}).Info("Player disconnect from lobby")

		// unregister client from lobby
		MainLobby.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {

		lobbyRequest := LobbyRequestMessage{}
		err := c.conn.ReadJSON(&lobbyRequest)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)

			}
			break
		}
		switch lobbyRequest.MessageType {
		case "JOIN_GAME":
			JoinGame(lobbyRequest, c)
		case "CREATE_GAME":
			CreateGame(c)
		case "SET_NAME":
			SetName(lobbyRequest, c)
		case "SEND_MSG":
			MainLobby.broadcast <- []byte(lobbyRequest.Message)
		}
	}
}

// gameClientRead reads messages from game connections
func (c *Client) gameClientRead() {
	defer func() {
		// unregister client from game
		log.WithFields(log.Fields{
			"playerID": c.ID,
		}).Info("Player disconnect from game")

		for _, g := range MainLobby.Games {
			g.unregister <- c
		}
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		gameRequest := GameRequestMessage{}
		err := c.conn.ReadJSON(&gameRequest)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)

			}
			break
		}
		switch gameRequest.MessageType {
		case "GAME_MOVE":
			GameMove(gameRequest, c)
		case "GAME_SEND_MESSAGE":
			MainLobby.broadcast <- []byte(gameRequest.Message)
		}

	}
}
