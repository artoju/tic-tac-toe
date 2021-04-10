package state

import (
	"context"
	"errors"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type RedisStateHandler struct {
	DB *redis.Client
}

var ctx = context.Background()

func (r RedisStateHandler) CreateState(req Game, c chan GameResponse) {
	val, err := r.DB.Get(ctx, req.ID).Result()
	if err != nil && err != redis.Nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error sending close message")
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	if val != "" {
		c <- GameResponse{
			Error: errors.New("Game already exists"),
		}
		return
	}

	str := req.Board + ":" + req.Status + ":" + req.NextPlayer
	err = r.DB.Set(ctx, req.ID, str, 0).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error sending close message")
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}

	c <- GameResponse{
		Error: nil,
		Game:  req,
	}
}
func (r RedisStateHandler) GetState(req Game, c chan GameResponse) {
	val, err := r.DB.Get(ctx, req.ID).Result()
	if err != nil {
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}

	game := Game{}
	parts := strings.Split(val, ":")
	game.Board, game.Status, game.NextPlayer = parts[0], parts[1], parts[2]
	game.ID = req.ID

	c <- GameResponse{
		Error: nil,
		Game:  game,
	}
}

func (r RedisStateHandler) SaveState(req Game, c chan GameResponse) {
	str := req.Board + ":" + req.Status + ":" + req.NextPlayer
	err := r.DB.Set(ctx, req.ID, str, 0).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error sending close message")
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	c <- GameResponse{
		Game:  req,
		Error: nil,
	}

}

func (r RedisStateHandler) DeleteState(id string, c chan GameResponse) {
	err := r.DB.Del(ctx, id).Err()
	if err != nil {
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	c <- GameResponse{
		Error: nil,
	}

}

func (r RedisStateHandler) GetStates() ([]Game, error) {
	games := make([]Game, 0)

	val, err := r.DB.Do(ctx, "SCAN", "0").Result()
	if err != nil {
		if err == redis.Nil {
			return games, err
		}
		panic(err)
	}
	for _, k := range val.([]interface{})[1].([]interface{}) {
		val, err := r.DB.Get(ctx, k.(string)).Result()
		if err != nil {
			return games, err
		}

		game := Game{}

		parts := strings.Split(val, ":")
		game.Board, game.Status = parts[0], parts[1]
		game.ID = k.(string)
		games = append(games, game)
	}

	return games, nil
}

func (r RedisStateHandler) ValidatePlayer(token string, id string) (*string, error) {
	val, err := r.DB.Get(ctx, "token:"+id).Result()
	if err != nil {
		if err == redis.Nil {
			UUID := uuid.Must(uuid.NewRandom()).String()

			err := r.DB.Set(ctx, "token:"+id, UUID, 0).Err()
			if err != nil {
				return nil, err
			}
			return &UUID, nil
		}
		return nil, err
	}
	if val != token {
		return nil, errors.New("Unauthorized")
	}
	return nil, nil
}
