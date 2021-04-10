package handlers

import (
	"encoding/json"
	"net/http"

	game "github.com/artoju/tic-tac-toe/game/actions"
	"github.com/artoju/tic-tac-toe/utils"
	"github.com/gorilla/mux"
)

func (f *StateHandler) DeleteGameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	gameID := vars["id"]
	token, err := utils.GetBearerToken(*r)
	_, err = f.GameState.ValidatePlayer(*token, gameID)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"reason": "Token unauthorized"})
		return
	}
	err = game.DeleteGame(gameID, f.GameState)
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
	json.NewEncoder(w).Encode(map[string]interface{}{"response": "success"})
}
