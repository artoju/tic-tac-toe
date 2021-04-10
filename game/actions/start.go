package game

import (
	"github.com/artoju/tic-tac-toe/game/state"
	uuid "github.com/google/uuid"
)

var GameChannel = make(chan state.GameResponse)

func CreateGame(sh state.GameState, isSinglePlayer bool) (*state.Game, error) {

	UUID := uuid.Must(uuid.NewRandom()).String()

	game := state.Game{
		ID:         UUID,
		Status:     "RUNNING",
		Board:      "---------",
		NextPlayer: "X",
	}

	go sh.SaveState(game, GameChannel)

	res := <-GameChannel
	if res.Error != nil {
		return nil, res.Error
	}

	return &res.Game, nil
}
