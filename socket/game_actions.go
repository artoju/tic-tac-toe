package socket

import gameAction "github.com/artoju/tic-tac-toe/game/actions"

// GameMove makes a move for client c.
func GameMove(gameRequest GameRequestMessage, c *Client) {

	// Check if both players present
	if len(c.game.Players) != 2 {
		errMessage := ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: "Game requires both players"}
		c.send <- errMessage
		return
	}
	req := c.game.gameObj
	req.Board = gameRequest.Message
	gameState, err := gameAction.UpdateGame(req, *MainLobby.GameStateHandler, false)
	if err != nil {
		c.send <- ErrorMessage{MessageType: "GAME_ERR_MESSAGE", Message: err.Error()}
		return
	}
	// Broadcast move to game.
	gameMoveMessage := GameMessage{
		Players:     []LobbyPlayer{},
		MessageType: "GAME_PLAYER_MOVED",
		Message:     "",
		NextPlayer:  gameState.NextPlayer,
		Board:       gameState.Board,
		GameStatus:  gameState.Status,
	}
	c.game.broadcast <- gameMoveMessage
}
