package main

import (
	"testing"

	"github.com/artoju/tic-tac-toe/config"
	"github.com/artoju/tic-tac-toe/game/state"
	"github.com/artoju/tic-tac-toe/redis"
)

func TestRedis(t *testing.T) {
	conf, err := config.Init()
	if err != nil {
		t.Log("Read config error " + err.Error())
		t.FailNow()
	}
	client, err := redis.Init(conf)
	if err != nil {
		t.Log("Init redis error " + err.Error())
		t.FailNow()
	}
	handler := state.RedisStateHandler{DB: client}
	channel := make(chan state.GameResponse, 0)
	req := state.Game{
		ID:     "12345",
		Board:  "---------",
		Status: "RUNNING",
	}

	go handler.CreateState(req, channel)
	res := <-channel
	if res.Error != nil {
		t.Log("Create state error " + res.Error.Error())
		t.FailNow()
	}
	go handler.CreateState(state.Game{ID: "123456", Board: "---------", Status: "RUNNING"}, channel)
	res = <-channel
	if res.Error != nil {
		t.Log("Create state error " + res.Error.Error())
		t.FailNow()
	}
	t.Log(res.Game.Board)

	req.Board = "-O-------"
	go handler.SaveState(req, channel)
	res = <-channel
	if res.Error != nil {
		t.Log("Delete state error " + res.Error.Error())
		t.FailNow()
	}

	go handler.GetState(req, channel)
	res = <-channel
	if res.Game.Board != "-O-------" {
		t.Log("Delete state error " + res.Error.Error())
		t.FailNow()
	}

	games, err := handler.GetStates()
	if err != nil {
		t.Log("Delete state error " + err.Error())
		t.FailNow()
	}
	if len(games) != 2 {
		t.Log("Expected length 2, got: ", len(games))
		t.FailNow()
	}

	go handler.DeleteState(req.ID, channel)
	res = <-channel
	if res.Error != nil {
		t.Log("Delete state error " + res.Error.Error())
		t.FailNow()
	}
	go handler.DeleteState("123456", channel)
	res = <-channel
	if res.Error != nil {
		t.Log("Delete state error " + res.Error.Error())
		t.FailNow()
	}

}
