// Package socket provides multiplayer functionality for game
package socket

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/artoju/tic-tac-toe/auth"
	gameAction "github.com/artoju/tic-tac-toe/game/actions"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// MainLobby is the single running instance of pre-game lobby
var MainLobby = NewLobby()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return strings.Contains(r.Header.Get("Origin"), "localhost")
	},
}

func NewLobby() *Lobby {
	return &Lobby{
		Games:          make([]*OnlineGame, 0),
		broadcast:      make(chan interface{}),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		registerGame:   make(chan *OnlineGame),
		unregisterGame: make(chan string),
	}
}

func NewGame() *OnlineGame {
	return &OnlineGame{
		broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		Players:    make(map[*Client]bool),
	}
}

// sendUpdateLobby broadcasts lobby update message.
func sendUpdateLobby() {
	updateLobby := LobbyMessage{Players: listPlayers(), MessageType: "UPDATE_LOBBY", Message: "", Games: listGames()}
	MainLobby.broadcast <- updateLobby
}

// Removes empty game after 10 seconds.
func sendRemoveGame(ID string) {
	time.Sleep(10 * time.Second)
	for _, g := range MainLobby.Games {
		if len(g.Players) == 0 {
			MainLobby.unregisterGame <- ID
		}
	}

}

// Run launches lobby loop.
func (h *Lobby) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				go sendUpdateLobby()
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
				}
			}
		case newGame := <-h.registerGame:
			h.Games = append(h.Games, newGame)

		case gameID := <-h.unregisterGame:
			var index int
			for idx, g := range h.Games {
				if g.ID == gameID {
					index = idx
				}
			}
			h.Games = append(h.Games[:index], h.Games[index+1:]...)

		}
	}
}

// run launches game loop.
func (game *OnlineGame) run() {
	for {
		select {
		case player := <-game.register:

			log.WithFields(log.Fields{
				"gameID":   game.ID,
				"playerID": player.ID,
			}).Info("Register new player")

			game.Players[player] = true

		case player := <-game.unregister:
			log.WithFields(log.Fields{
				"gameID":   game.ID,
				"playerID": player.ID,
			}).Info("Game unregister player")

			var p *Client

			for c := range game.Players {
				if c.ID == player.ID {
					p = c
				}
			}
			if _, ok := game.Players[p]; ok {
				delete(game.Players, p)
				if len(game.Players) == 0 {
					go sendRemoveGame(game.ID)
				}
			}
		case message := <-game.broadcast:
			for player := range game.Players {
				select {
				case player.send <- message:
				default:
					close(player.send)
					delete(game.Players, player)
				}
			}

		}
	}
}

// Lists active games.
func listGames() []LobbyGame {
	games := make([]LobbyGame, 0)

	for _, game := range MainLobby.Games {
		players := make([]LobbyPlayer, 0)
		for c := range game.Players {
			players = append(players, LobbyPlayer{ID: c.ID})
		}
		games = append(games, LobbyGame{PlayerCount: len(game.Players), GameID: game.ID, Players: players})
	}
	return games
}

// Lists active lobby clients.
func listPlayers() []LobbyPlayer {
	players := make([]LobbyPlayer, 0)
	for player := range MainLobby.clients {
		players = append(players, LobbyPlayer{ID: player.ID, Name: player.Name})
	}
	return players
}

// ShortID generates short id's for games and players.
func ShortID() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

// LobbyHandler serves lobby websocket
func LobbyHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"context": "LobbyHandler websocket upgrade",
			"error":   err.Error(),
		}).Error("LobbyHandler error")
		return
	}
	players := make([]LobbyPlayer, 0)
	for player := range MainLobby.clients {
		players = append(players, LobbyPlayer{ID: player.ID, Name: player.Name})
	}
	id := ShortID()

	log.WithFields(log.Fields{
		"playerID": id,
	}).Info("New lobby user")

	newPlayerMessage := LobbyMessage{MessageType: "LOBBY_NEW_PLAYER", Message: id, Players: listPlayers()}

	MainLobby.broadcast <- newPlayerMessage
	client := &Client{conn: conn, send: make(chan interface{}, 256), ID: id}

	joinedLobbyMessage := LobbyMessage{Players: players, MessageType: "LOBBY_JOIN", Message: id, Games: listGames()}

	client.send <- joinedLobbyMessage
	MainLobby.register <- client

	go client.clientWrite()
	go client.lobbyClientRead()
}

// CloseWithMessage closes connection with a message.
func CloseWithMessage(c *websocket.Conn, msg string) {
	cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, msg)
	if err := c.WriteMessage(websocket.CloseMessage, cm); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error sending close message")
	}
	c.Close()
}

// GameHandler serves game websocket connection.
func GameHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"context": "Gamehandler websocket upgrade",
			"error":   err.Error(),
		}).Error("Gamehandler error")
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	token := vars["token"]
	var g *OnlineGame

	if id == "" || token == "" {
		CloseWithMessage(conn, "No game ID or token specified")
		return
	}

	for _, game := range MainLobby.Games {
		if game.ID == id && len(game.Players) != 2 {
			g = game
		}
	}
	if g == nil {
		CloseWithMessage(conn, "No game found for ID")
		return
	}

	playerSign, err := auth.AuthenticatePlayer(*MainLobby.GameStateHandler, token, g.gameObj.ID)
	if err != nil {
		CloseWithMessage(conn, "Player not authorized for game")
		return
	}

	playerId := ShortID()
	client := &Client{conn: conn, send: make(chan interface{}, 256), ID: playerId}

	g.register <- client
	client.game = g

	log.WithFields(log.Fields{
		"playerID": playerId,
		"gameID":   g.ID,
	}).Info("Player joined game")

	gameState, err := gameAction.GetGame(g.gameObj, *MainLobby.GameStateHandler)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error finding gmae")
		CloseWithMessage(conn, "No game found for ID")
		return
	}

	gameJoinMessage := GameJoinMessage{
		MessageType: "GAME_PLAYER_JOINED",
		Board:       gameState.Board,
		PlayerSign:  *playerSign,
		NextPlayer:  gameState.NextPlayer,
		GameStatus:  gameState.Status,
	}

	client.send <- gameJoinMessage

	go client.clientWrite()
	go client.gameClientRead()
}
