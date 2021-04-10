package game

import "github.com/artoju/tic-tac-toe/game/state"

func GetGames(sh state.GameState) ([]state.Game, error) {
	return sh.GetStates()
}
