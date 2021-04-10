package handlers

import (
	"encoding/json"
	"net/http"

	game "github.com/artoju/tic-tac-toe/game/actions"
	"github.com/artoju/tic-tac-toe/game/state"
	"github.com/artoju/tic-tac-toe/utils"
	"github.com/gorilla/mux"
)

func (f *StateHandler) UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	var req state.Game

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	gameID := vars["id"]
	if gameID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": "ID is missing from path"})
		return
	}

	token, err := utils.GetBearerToken(*r)
	_, err = f.GameState.ValidatePlayer(*token, gameID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": "Token unauthorized"})
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": err.Error()})
		return
	}

	moveRequest := state.Game{
		ID:    gameID,
		Board: req.Board,
	}

	game, err := game.UpdateGame(moveRequest, f.GameState, f.IsSinglePlayer)
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

	json.NewEncoder(w).Encode(game)
}
