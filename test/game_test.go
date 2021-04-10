package main

import (
	"os"
	"testing"

	game "github.com/artoju/tic-tac-toe/game/actions"
	minimax "github.com/artoju/tic-tac-toe/game/minimax"
	"github.com/artoju/tic-tac-toe/game/state"
)

var dirname, _ = os.Getwd()

var f = state.FileStateHandler{Filepath: dirname + "/games/"}

func TestFindBestMove(t *testing.T) {
	board := "O---X----"
	move := minimax.FindBestMove(board, "X", "O")
	if move != "OX--X----" {
		t.Log("Expected move OX--X----, got: " + move)
		t.FailNow()
	}
	t.FailNow()

}

func TestConcurrency(t *testing.T) {
	gameChannel := make(chan state.GameResponse)
	req := state.Game{
		ID:     "",
		Board:  "---------",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	g2, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	req11 := state.Game{
		ID:     g2.ID,
		Board:  "----X----",
		Status: "",
	}

	req22 := state.Game{
		ID:     g2.ID,
		Board:  "---X-----",
		Status: "",
	}
	req33 := state.Game{
		ID:     g2.ID,
		Board:  "--------X",
		Status: "",
	}
	req44 := state.Game{
		ID:     g2.ID,
		Board:  "X---O---X",
		Status: "",
	}
	req.ID = g.ID

	res, err := game.GetGame(req, f)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}
	req.ID = res.ID
	req.Board = "----X----"

	req2 := state.Game{
		ID:     res.ID,
		Board:  "---X-----",
		Status: "",
	}
	req3 := state.Game{
		ID:     res.ID,
		Board:  "--------X",
		Status: "",
	}
	req4 := state.Game{
		ID:     res.ID,
		Board:  "X---O---X",
		Status: "",
	}
	go func() {
		g, err := game.UpdateGame(req, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req2, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req4, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req3, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req11, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req22, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req44, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()
	go func() {
		g, err := game.UpdateGame(req33, f, true)
		if g == nil {
			g = &state.Game{}
		}
		gameChannel <- state.GameResponse{
			Game:  *g,
			Error: err,
		}
	}()

	counter := 0
	for {
		select {
		case res := <-gameChannel:
			counter++
			if res.Error != nil {
				t.Log(res.Error)
			} else {
				t.Log(res.Game.Board)
			}
			if counter == 8 {
				t.FailNow()
				break
			}

		}
	}

}

func TestStartGame(t *testing.T) {

	state, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	err = game.DeleteGame(state.ID, f)
	if err != nil {
		t.FailNow()
	}

}
func TestListGames(t *testing.T) {

	_, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	_, err = game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	games, err := game.GetGames(f)
	if len(games) != 2 {
		t.Logf("Number of games found: %d", len(games))
		t.FailNow()
	}
	for _, g := range games {
		err = game.DeleteGame(g.ID, f)
		if err != nil {
			t.FailNow()
		}
	}

}
func TestDeleteGame(t *testing.T) {
	req := state.Game{
		ID:     "",
		Board:  "---------",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}
	req.ID = g.ID
	err = game.DeleteGame(req.ID, f)
	if err != nil {
		t.FailNow()
	}
}

func TestStartGameWithMoveO(t *testing.T) {
	req := state.Game{
		ID:     "",
		Board:  "O--------",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}
	req.ID = g.ID
	err = game.DeleteGame(req.ID, f)
	if err != nil {
		t.FailNow()
	}
}
func TestStartGameWithMoveX(t *testing.T) {
	req := state.Game{
		ID:     "",
		Board:  "----X----",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}
	req.ID = g.ID
	err = game.DeleteGame(req.ID, f)
	if err != nil {
		t.FailNow()
	}
}

func TestMoveGame(t *testing.T) {
	req := state.Game{
		ID:     "",
		Board:  "---------",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	req.ID = g.ID

	res, err := game.GetGame(req, f)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}
	req.ID = res.ID
	req.Board = "----X----"
	state, err := game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}
	if state.Board != "O---X----" {
		t.Log("Expected board: O---X----, actual: " + state.Board)
		t.FailNow()
	}
	err = game.DeleteGame(req.ID, f)
	if err != nil {
		t.Log("Deleting game error: " + err.Error())
		t.FailNow()
	}

}

func TestLostGame(t *testing.T) {
	req := state.Game{
		ID:     "",
		Board:  "---------",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	req.ID = g.ID
	req.Board = "----X----"
	g, err = game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}
	t.Log(g.Board)
	req.Board = "O---X---X"
	state, err := game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}
	if state.Board != "O-O-X---X" {
		t.Log("Expected board: O-O-X---X, actual: " + state.Board)
		t.FailNow()
	}
	req.Board = "O-O-X--XX"
	state, err = game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}

	if state.Status != "O_WON" {
		t.Log("Expected status: O_WON, actual: " + state.Status)
		t.FailNow()
	}
	req.Board = "OOO-X-XXX"

	state, err = game.UpdateGame(req, f, true)
	if err == nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}
	err = game.DeleteGame(req.ID, f)
	if err != nil {
		t.Log("Deleting game error: " + err.Error())
		t.FailNow()
	}

}

func TestDrawGame(t *testing.T) {
	req := state.Game{
		ID:     "",
		Board:  "---------",
		Status: "",
	}

	g, err := game.CreateGame(f, true)
	if err != nil {
		t.Log("Create game error: " + err.Error())
		t.FailNow()
	}

	req.ID = g.ID

	req.Board = "----X----"
	g, err = game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}

	req.Board = "O---X---X"
	state, err := game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}
	if state.Board != "O-O-X---X" {
		t.Log("Expected board: O-O-X---X, actual: " + state.Board)
		t.FailNow()
	}
	req.Board = "OXO-X---X"
	state, err = game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}

	if state.Board != "OXO-X--OX" {
		t.Log("Expected board: OXO-X--OX, actual: " + state.Board)

	}
	req.Board = "OXOXX--OX"
	state, err = game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}

	if state.Board != "OXOXXO-OX" {
		t.Log("Expected board: OXOXXO-OX, actual: " + state.Board)
		t.FailNow()
	}

	req.Board = "OXOXXOXOX"
	state, err = game.UpdateGame(req, f, true)
	if err != nil {
		t.Log("Updating game error: " + err.Error())
		t.FailNow()
	}
	if state.Status != "DRAW" {
		t.Log("Status is not DRAW, actual status: " + state.Status)
		t.FailNow()
	}
	if state.Board != "OXOXXOXOX" {
		t.Log("Expected board: OXOXXOXOX, actual: " + state.Board)
		t.FailNow()
	}
	_, err = game.UpdateGame(req, f, true)
	if err == nil {
		t.Log("Expected error updating")
		t.FailNow()
	}

	err = game.DeleteGame(req.ID, f)
	if err != nil {
		t.Log("Error deleting: " + err.Error())
		t.FailNow()
	}
}
