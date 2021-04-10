package game

import (
	"errors"
)

func ValidateFirstMove(board string, move string) (string, error) {
	if len(move) != 9 {
		return "", errors.New("Board must be length of 9")
	}
	boardRunes := []rune(board)
	moveRunes := []rune(move)
	foundChange := false
	player := ""
	Xcount := 0
	Ocount := 0
	for i := 0; i < 9; i++ {
		boardChar := string(boardRunes[i])
		moveChar := string(moveRunes[i])
		if boardChar == "O" {
			Ocount++
		}
		if boardChar == "X" {
			Xcount++
		}
		if boardChar != moveChar {
			if boardChar != "-" {
				return "", errors.New("Invalid move")
			}
			if foundChange {
				return "", errors.New("Found more than one change in board")
			}
			foundChange = true
			player = moveChar
		}

	}
	delta := Ocount - Xcount

	if delta < -1 || delta > 1 {
		return "", errors.New("More than one move found")
	}
	if !foundChange {
		return "", errors.New("Found no moves")
	}
	return player, nil
}

func ValidateMove(board string, move string, mover string) error {
	if len(move) != 9 {
		return errors.New("Board must be length of 9")
	}
	boardRunes := []rune(board)
	moveRunes := []rune(move)
	foundChange := false
	Xcount := 0
	Ocount := 0
	for i := 0; i < 9; i++ {
		boardChar := string(boardRunes[i])
		moveChar := string(moveRunes[i])
		if boardChar == "O" {
			Ocount++
		}
		if boardChar == "X" {
			Xcount++
		}
		if boardChar != moveChar {
			if boardChar != "-" {
				return errors.New("Invalid move")
			}
			if foundChange {
				return errors.New("Found more than one change in board")
			}
			if moveChar != mover {
				return errors.New("Invalid move: expected " + mover + ", got " + moveChar)
			}
			foundChange = true
		}

	}
	if !foundChange {
		return errors.New("Found no moves")
	}
	return nil
}
