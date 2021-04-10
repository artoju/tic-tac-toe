package socket

import (
	"github.com/artoju/tic-tac-toe/game/state"
	"github.com/gorilla/websocket"
)

// Client represents lobby or game connection.
type Client struct {
	// Identificates client.
	ID string

	// Access client's current game. Nil if empty.
	game *OnlineGame

	// Client socket connection.
	conn *websocket.Conn

	// Send private messages to client.
	send chan interface{}
}

// Lobby represents the pregame lobby and holds client connections
// and channels for broadcasting messages and creating games.
type Lobby struct {
	// Registered games.
	Games []*OnlineGame

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan interface{}

	// Register lobby clients.
	register chan *Client

	// Unregister lobby clients.
	unregister chan *Client

	// Register new initialized game.
	registerGame chan *OnlineGame

	// Remove game from game selection.
	unregisterGame chan string

	// Handles game state management.
	GameStateHandler *state.GameState
}

// OnlineGame holds the player client connections
// and channels for broadcasting game actions.
type OnlineGame struct {
	// Identificates OnlineGame.
	ID string

	// Connected player clients.
	Players map[*Client]bool

	// Actual game state.
	gameObj state.Game

	// Send struct type messages to be broadcasted as JSON.
	broadcast chan interface{}

	// Register player to game.
	register chan *Client

	// Unregister player from game.
	unregister chan *Client
}

// Outgoing messages to lobby clients.
type LobbyMessage struct {
	MessageType string      `json:"messageType"`
	Message     string      `json:"message"`
	Players     []string    `json:"players"`
	Games       []LobbyGame `json:"games"`
}

// Outgoing response message for joining game.
type GameJoinMessage struct {
	MessageType string `json:"messageType"`
	PlayerSign  string `json:"playerSign"`
	Board       string `json:"board"`
	NextPlayer  string `json:"nextPlayer"`
	GameStatus  string `json:"gameStatus"`
}

// Outgoing game action broadcast message.
type GameMessage struct {
	MessageType string   `json:"messageType"`
	Message     string   `json:"message"`
	Players     []string `json:"players"`
	Board       string   `json:"board"`
	NextPlayer  string   `json:"nextPlayer"`
	GameStatus  string   `json:"gameStatus"`
}

// Game's lobby presentation.
type LobbyGame struct {
	PlayerCount int           `json:"playerCount"`
	GameID      string        `json:"gameId"`
	Players     []LobbyPlayer `json:"players"`
}

// LobbyGame's player presentation.
type LobbyPlayer struct {
	ID string `json:"id"`
}

// Lobby client request
type LobbyRequestMessage struct {
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
}

// Incoming request message for game action.
type GameRequestMessage struct {
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
}

// Generic outoing error response message.
type ErrorMessage struct {
	MessageType string `json:"messageType"`
	Message     string `json:"message"`
}
