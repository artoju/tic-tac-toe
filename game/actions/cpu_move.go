package game

import (
	"github.com/artoju/tic-tac-toe/game/minimax"
	"github.com/artoju/tic-tac-toe/game/state"
)

// CpuMove makes a move for the player that is set to be next.
func CpuMove(req state.Game, sh state.GameState) (*state.Game, error) {
	go sh.GetState(req, GameChannel)

	res := <-GameChannel
	if res.Error != nil {
		return nil, res.Error
	}
	game := res.Game
	board := game.Board
	status := game.Status
	cpu := game.NextPlayer

	player := "X"
	if cpu == "X" {
		player = "O"
	}

	board = minimax.FindBestMove(req.Board, cpu, player)
	result := minimax.Evaluate(minimax.ConvertStrToBoard(board), cpu, player)
	if result == -10 {
		status = player + "_WON"
	} else if result == 10 {
		status = cpu + "_WON"
	} else if !minimax.IsMovesLeft(minimax.ConvertStrToBoard(board)) {
		status = "DRAW"
	}
	req.Board = board
	req.Status = status
	req.NextPlayer = player

	go sh.SaveState(req, GameChannel)
	res = <-GameChannel
	if res.Error != nil {
		return nil, res.Error
	}

	return &res.Game, nil
}
