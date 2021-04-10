package game

import "github.com/artoju/tic-tac-toe/game/state"

func DeleteGame(id string, sh state.GameState) error {
	go sh.DeleteState(id, GameChannel)

	res := <-GameChannel
	if res.Error != nil {
		return res.Error
	}
	return nil
}
