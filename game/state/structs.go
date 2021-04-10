package state

// GameState wraps necessary game state actions.
type GameState interface {
	CreateState(req Game, c chan GameResponse)
	GetState(req Game, c chan GameResponse)
	SaveState(req Game, c chan GameResponse)
	DeleteState(id string, c chan GameResponse)
	GetStates() ([]Game, error)
	ValidatePlayer(token string, id string) (*string, error)
}

// Game acts as
type Game struct {
	// Identification as is set in selected
	// state saving method.
	ID string

	// Board represented as string e.g. "---------"
	Board string

	/*
	 Game status.
	 RUNNING - ongoing game
	 O_WON - game over, won by O
	 X_WON - game over, won by X
	 DRAW - game over, draw
	*/
	Status string

	// Player allowed to make next move noted
	// by player sign O or X.
	NextPlayer string
}

// GameResponse is the response received from committing
// GameState actions.
type GameResponse struct {
	Game  Game
	Error error
}
