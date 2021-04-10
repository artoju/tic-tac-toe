package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/artoju/tic-tac-toe/auth"
	game "github.com/artoju/tic-tac-toe/game/actions"
)

func (f *StateHandler) StartGameHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	g, err := game.CreateGame(f.GameState, f.IsSinglePlayer)
	if err != nil {
		if err.Error() == "Internal error" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{"reason": err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": err.Error()})
		return
	}
	token, err := auth.CreateToken(f.GameState, g.ID, g.NextPlayer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"ID": g.ID, "token": token})
}
