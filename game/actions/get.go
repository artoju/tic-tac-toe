package game

import "github.com/artoju/tic-tac-toe/game/state"

func GetGame(req state.Game, sh state.GameState) (*state.Game, error) {
	res := state.GameResponse{}

	go sh.GetState(req, GameChannel)

	res = <-GameChannel
	if res.Error != nil {
		return nil, res.Error
	}
	return &res.Game, nil
}
