package main

import (
	"testing"

	minimax "github.com/artoju/tic-tac-toe/game/minimax"
)

func TestEvaluateDynamic(t *testing.T) {
	res := minimax.EvaluateDynamic(
		[][]string{
			{"X", "X", "O"},
			{"X", "O", "X"},
			{"O", "X", "X"},
		},
		"O", "X",
	)
	if res != 10 {
		t.FailNow()
	}
	res = minimax.EvaluateDynamic(
		[][]string{
			{"X", "X", "X", "O"},
			{"X", "X", "O", "X"},
			{"O", "O", "O", "X"},
			{"O", "O", "X", "X"},
		},
		"O", "X",
	)
	if res != 10 {
		t.FailNow()
	}
	t.Log("result", res)
}
