package game

import (
	"errors"

	"github.com/EagleChen/mapmutex"
	"github.com/artoju/tic-tac-toe/game/minimax"
	"github.com/artoju/tic-tac-toe/game/state"
)

var mutex = mapmutex.NewMapMutex()

// UpdateGame updates game synchronically based on the game ID.
func UpdateGame(req state.Game, sh state.GameState, isSinglePlayer bool) (*state.Game, error) {

	if mutex.TryLock(req.ID) {
		defer mutex.Unlock(req.ID)
		go sh.GetState(req, GameChannel)

		res := <-GameChannel
		if res.Error != nil {
			return nil, res.Error
		}
		game := res.Game
		board := game.Board
		status := game.Status
		player := game.NextPlayer
		if status != "RUNNING" {
			return nil, errors.New("Game is ended")
		}

		err := ValidateMove(board, req.Board, player)
		if err != nil {
			return nil, err
		}
		game.Board = req.Board
		opponent := "X"
		if player == "X" {
			opponent = "O"
		}
		game.NextPlayer = opponent

		result := minimax.EvaluateDynamic(minimax.ConvertStrToBoard(game.Board), player, opponent)
		if result == 10 {
			status = player + "_WON"
		} else if result == -10 {
			status = opponent + "_WON"
		} else {
			if minimax.IsMovesLeft(minimax.ConvertStrToBoard(game.Board)) {
				status = "RUNNING"
			} else {
				status = "DRAW"
			}
		}
		game.Status = status
		if isSinglePlayer {
			go sh.SaveState(game, GameChannel)
			res = <-GameChannel
			if res.Error != nil {
				return nil, res.Error
			}

			return CpuMove(res.Game, sh)
		}

		go sh.SaveState(game, GameChannel)
		res = <-GameChannel
		if res.Error != nil {
			return nil, res.Error
		}

		return &res.Game, nil
	}
	return &state.Game{}, nil
}
