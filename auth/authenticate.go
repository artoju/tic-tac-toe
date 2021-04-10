package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/artoju/tic-tac-toe/game/state"
	"github.com/google/uuid"
)

var ctx = context.Background()

//AuthenticatePlayer authenticates player by attempting to match
//token with gameId in selected gameState handler. Returns error
//if match is not made or game is not found. If player is authenticated
//player sign is returned.
func AuthenticatePlayer(gameState state.GameState, token string, gameId string) (*string, error) {
	stateHandlerType := fmt.Sprintf("%T", gameState)
	switch stateHandlerType {
	case "state.RedisStateHandler":
		return RedisAuthenticationHandler(gameState.(state.RedisStateHandler), token, gameId)
	case "state.FileStateHandler":
		return nil, errors.New("Authentication for state handler not yet implemented")
	case "state.DBStateHandler":
		return nil, errors.New("Authentication for state handler not yet implemented")
	default:
		return nil, errors.New("Found no state handler")
	}
}

//CreateToken creates authentication token for player for game with gameId
//via gameState handler. Player's sign is saved in playerSign.
func CreateToken(gameState state.GameState, gameId string, playerSign string) (*string, error) {
	stateHandlerType := fmt.Sprintf("%T", gameState)
	switch stateHandlerType {
	case "state.RedisStateHandler":
		return RedisCreateTokenHandler(gameState.(state.RedisStateHandler), gameId, playerSign)
	case "state.FileStateHandler":
		return nil, errors.New("Token creation for state handler not yet implemented")
	case "state.DBStateHandler":
		return nil, errors.New("Token creation for state handler not yet implemented")
	default:
		return nil, errors.New("Found no state handler")
	}
}

//DeleteToken removes authentication token from gameState handler.
func DeleteToken(gameState state.GameState, token string) error {
	stateHandlerType := fmt.Sprintf("%T", gameState)
	switch stateHandlerType {
	case "state.RedisStateHandler":
		return RedisDeleteTokenHandler(gameState.(state.RedisStateHandler), token)
	case "state.FileStateHandler":
		return errors.New("Token deletion for state handler not yet implemented")
	case "state.DBStateHandler":
		return errors.New("Token deletion for state handler not yet implemented")
	default:
		return errors.New("Found no state handler")
	}
}

//RedisAuthenticationHandler is redis specific authentication handler.
func RedisAuthenticationHandler(r state.RedisStateHandler, token string, id string) (*string, error) {
	val, err := r.DB.Get(ctx, "token:"+token).Result()
	if err != nil {
		return nil, err
	}
	parts := strings.Split(val, ":")

	if parts[0] != id {
		return nil, errors.New("Unauthorized")
	}
	return &parts[1], nil
}

//RedisCreateTokenHandler is redis specific handler for token creation.
func RedisCreateTokenHandler(r state.RedisStateHandler, gameId string, playerSign string) (*string, error) {

	UUID := uuid.Must(uuid.NewRandom()).String()

	err := r.DB.Set(ctx, "token:"+UUID, gameId+":"+playerSign, 0).Err()
	if err != nil {
		return nil, err
	}
	return &UUID, nil

}

//RedisDeleteTokenHandler is redis specific handler for removing tokens.
func RedisDeleteTokenHandler(r state.RedisStateHandler, token string) error {
	err := r.DB.Del(ctx, "token:"+token).Err()
	if err != nil {
		return err
	}

	return nil
}
