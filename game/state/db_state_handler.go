package state

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DBStateHandler struct {
	DB *dynamodb.DynamoDB
}

type GameDBInput struct {
	ID        string
	Board     string
	Status    string
	UpdatedAt string
	Player    string
}

func (dbh DBStateHandler) CreateState(req Game, c chan GameResponse) {
	gameInput := gameToDBInput(req)
	newGame, err := dynamodbattribute.MarshalMap(gameInput)
	if err != nil {
		c <- GameResponse{
			Error: errors.New("No game found for ID"),
		}
		return
	}
	input := &dynamodb.PutItemInput{
		Item:      newGame,
		TableName: aws.String("tic-tac-toe"),
	}
	_, err = dbh.DB.PutItem(input)
	if err != nil {
		c <- GameResponse{
			Error: errors.New("No game found for ID"),
		}
		return
	}
}
func (dbh DBStateHandler) GetState(req Game, c chan GameResponse) {
	gameInput := gameToDBInput(req)
	newGame, err := dynamodbattribute.MarshalMap(gameInput)
	if err != nil {
		c <- GameResponse{
			Error: errors.New("No game found for ID"),
		}
		return
	}
	input := &dynamodb.GetItemInput{
		Key:       newGame,
		TableName: aws.String("tic-tac-toe"),
	}
	output, err := dbh.DB.GetItem(input)
	if err != nil {
		c <- GameResponse{
			Error: errors.New("No game found for ID"),
		}
		return
	}
	gameOutput := GameDBInput{}
	err = dynamodbattribute.UnmarshalMap(output.Item, &gameOutput)
	if err != nil {
		c <- GameResponse{
			Error: errors.New("No game found for ID"),
		}
		return
	}

	result := gameFromDBInput(gameOutput)
	c <- GameResponse{
		Error: nil,
		Game:  result,
	}
}

func (dbh DBStateHandler) SaveState(req Game, c chan GameResponse) {
	// TODO
}

func (dbh DBStateHandler) DeleteState(id string, c chan GameResponse) {
	// TODO
}

func (dbh DBStateHandler) GetStates() ([]Game, error) {
	games := make([]Game, 0)
	return games, nil
}

func gameToDBInput(game Game) GameDBInput {
	input := GameDBInput{}
	t := time.Now()
	ts := t.Format("20060102150405")

	input.UpdatedAt = ts
	input.Player = game.NextPlayer
	input.Board = game.Board
	input.ID = game.ID
	input.Status = game.Status

	return input

}

func gameFromDBInput(input GameDBInput) Game {
	game := Game{}
	game.NextPlayer = input.Player
	game.Board = input.Board
	game.ID = input.ID
	game.Status = input.Status

	return game
}

func (dbh DBStateHandler) ValidatePlayer(token string, id string) (*string, error) {
	return nil, nil
}
