package handlers

import (
	"encoding/json"
	"net/http"

	game "github.com/artoju/tic-tac-toe/game/actions"
)

func (f *StateHandler) GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	games, err := game.GetGames(f.GameState)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(games)
}
