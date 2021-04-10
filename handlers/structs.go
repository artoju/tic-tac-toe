// Package handlers contains API endpoint handlers for singleplayer.
package handlers

import "github.com/artoju/tic-tac-toe/game/state"

type StateHandler struct {
	GameState      state.GameState
	IsSinglePlayer bool
}
