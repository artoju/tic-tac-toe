package socket

import (
	"regexp"
	"time"

	"github.com/artoju/tic-tac-toe/auth"
	gameAction "github.com/artoju/tic-tac-toe/game/actions"
	log "github.com/sirupsen/logrus"
)

// JoinGame adds client c to the requested game if possible and
// sends join message accordingly.
func JoinGame(lobbyRequest LobbyRequestMessage, c *Client) {
	id := lobbyRequest.Message
	var game *OnlineGame
	for _, g := range MainLobby.Games {
		if g.ID == id {
			game = g
		}
	}
	// Check is game found
	if game == nil {
		errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game not found"}
		c.send <- errMessage
		return
	}

	// Check if game is full
	if len(game.Players) == 2 {
		errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game is full"}
		c.send <- errMessage
		return
	}

	// Check if game exists
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
		return
	}

	// Check if game is started
	if gameState.Board != "---------" {
		errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game already started"}
		c.send <- errMessage
		return
	}

	// Checks passed, player may be assigned to the game
	token, err := auth.CreateToken(*MainLobby.GameStateHandler, gameState.ID, "O")
	if err != nil {
		log.WithFields(log.Fields{
			"playerID": c.ID,
			"gameID":   gameState.ID,
			"error":    err.Error(),
		}).Error("Create token error")
		errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Internal error: token"}
		c.send <- errMessage
		return
	}

	// Unregister client
	MainLobby.unregister <- c
	updateLobbyMessage := LobbyMessage{Players: listPlayers(), MessageType: "UPDATE_LOBBY", Message: "", Games: listGames()}
	MainLobby.broadcast <- updateLobbyMessage

	// Send message with game id and token for joining
	gameJoinMessage := GameMessage{
		Players:     listPlayers(),
		MessageType: "LOBBY_JOINED_GAME",
		Message:     game.ID + ":" + *token,
		NextPlayer:  gameState.NextPlayer,
		Board:       gameState.Board,
		GameStatus:  gameState.Status,
	}
	c.send <- gameJoinMessage
	return

}

// CreateGame creates a new game and sends confirmation message to the client c.
func CreateGame(c *Client) {

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
		return
	}

	token, err := auth.CreateToken(*MainLobby.GameStateHandler, gameState.ID, "X")
	if err != nil {
		log.WithFields(log.Fields{
			"playerID": c.ID,
			"error":    err.Error(),
		}).Error("Create gametoken error")
		errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Unable to create game"}
		c.send <- errMessage
		return
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
}

// SetName sets c client's name and sends update message.
func SetName(lobbyRequest LobbyRequestMessage, c *Client) {
	re := regexp.MustCompile("[\t\n\f\r ]")
	name := re.ReplaceAllString(lobbyRequest.Message, "")
	c.Name = name
	updateNameMessage := UpdateNameMessage{Player: LobbyPlayer{Name: name, ID: c.ID}, MessageType: "UPDATE_PLAYER_NAME"}
	MainLobby.broadcast <- updateNameMessage
}

// SendLobbyChatMessage broadcasts a chat message to lobby.
func SendLobbyChatMessage(lobbyRequest LobbyRequestMessage, c *Client) {
	t := time.Now()
	ts := t.Format("15:04")
	re := regexp.MustCompile("[\t\n]")
	message := re.ReplaceAllString(lobbyRequest.Message, " ")
	chatMessage := ChatMessage{MessageType: "CHAT_MESSAGE", Message: message, Timestamp: ts, Sender: LobbyPlayer{ID: c.ID, Name: c.Name}}
	MainLobby.broadcast <- chatMessage
}
