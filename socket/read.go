package socket

import (
	"time"

	"github.com/artoju/tic-tac-toe/auth"
	gameAction "github.com/artoju/tic-tac-toe/game/actions"
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
		if lobbyRequest.MessageType == "JOIN_GAME" {
			id := lobbyRequest.Message
			var game *OnlineGame
			for _, g := range MainLobby.Games {
				if g.ID == id {
					game = g
				}
			}
			if game != nil && len(game.Players) != 2 {

				//Check if game is started
				gameState, err := gameAction.GetGame(game.gameObj, *MainLobby.GameStateHandler)
				if err != nil {
					log.WithFields(log.Fields{
						"context":  "Get gamestate",
						"playerID": c.ID,
						"gameID":   gameState.ID,
						"error":    err.Error(),
					}).Error("Join game error")
					errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Internal error: game"}
					c.send <- errMessage
					continue
				}

				if gameState.Board == "---------" {
					token, err := auth.CreateToken(*MainLobby.GameStateHandler, gameState.ID, "O")
					if err != nil {
						log.WithFields(log.Fields{
							"playerID": c.ID,
							"gameID":   gameState.ID,
							"error":    err.Error(),
						}).Error("Create token error")
						errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Internal error: token"}
						c.send <- errMessage
						continue
					}

					MainLobby.unregister <- c
					updateLobbyMessage := LobbyMessage{Players: listPlayers(), MessageType: "UPDATE_LOBBY", Message: "", Games: listGames()}
					MainLobby.broadcast <- updateLobbyMessage

					gameJoinMessage := GameMessage{
						Players:     listPlayers(),
						MessageType: "LOBBY_JOINED_GAME",
						Message:     game.ID + ":" + *token,
						NextPlayer:  gameState.NextPlayer,
						Board:       gameState.Board,
						GameStatus:  gameState.Status,
					}
					c.send <- gameJoinMessage
					continue

				} else {
					errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game already started"}
					c.send <- errMessage
					continue

				}

			} else {
				errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game is full"}
				c.send <- errMessage
				continue
			}

		} else if lobbyRequest.MessageType == "CREATE_GAME" {

			newGame := NewGame()
			gameID := ShortID()
			newGame.ID = gameID

			gameState, err := gameAction.CreateGame(*MainLobby.GameStateHandler, false)
			if err != nil {
				log.WithFields(log.Fields{
					"playerID": c.ID,
					"error":    err.Error(),
				}).Error("Create game error")
				errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Unable to create game"}
				c.send <- errMessage
				continue
			}

			token, err := auth.CreateToken(*MainLobby.GameStateHandler, gameState.ID, "X")
			if err != nil {
				log.WithFields(log.Fields{
					"playerID": c.ID,
					"error":    err.Error(),
				}).Error("Create gametoken error")
				errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Unable to create game"}
				c.send <- errMessage
				continue
			}

			newGame.gameObj = *gameState
			MainLobby.registerGame <- newGame
			go newGame.run()

			gameJoinMessage := GameMessage{
				Players:     listPlayers(),
				MessageType: "LOBBY_JOINED_GAME",
				Message:     gameID + ":" + *token,
			}
			c.send <- gameJoinMessage
		} else if lobbyRequest.MessageType == "SEND_MSG" {
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
		if gameRequest.MessageType == "GAME_MOVE" {
			if len(c.game.Players) != 2 {
				errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game requires both players"}
				c.send <- errMessage
				continue
			}
			req := c.game.gameObj
			req.Board = gameRequest.Message
			gameState, err := gameAction.UpdateGame(req, *MainLobby.GameStateHandler, false)
			if err != nil {
				c.send <- ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: err.Error()}
				continue
			}
			gameMoveMessage := GameMessage{
				Players:     []string{},
				MessageType: "GAME_PLAYER_MOVED",
				Message:     "",
				NextPlayer:  gameState.NextPlayer,
				Board:       gameState.Board,
				GameStatus:  gameState.Status,
			}
			c.game.broadcast <- gameMoveMessage
		} else if gameRequest.MessageType == "GAME_SEND_MESSAGE" {
			c.game.broadcast <- []byte(gameRequest.Message)
		}
	}
}
