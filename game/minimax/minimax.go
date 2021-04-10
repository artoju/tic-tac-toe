package minimax

// IsMovesLeft checks if board has any moves left.
func IsMovesLeft(board [][]string) bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "-" {
				return true
			}
		}
	}
	return false
}

// Evaluate returns board state from the view of player. -10 = loss / 10 = win / 0 = draw
func Evaluate(board [][]string, player string, opponent string) int {
	for row := 0; row < 3; row++ {
		if board[row][0] == board[row][1] &&
			board[row][1] == board[row][2] {
			if board[row][0] == player {
				return 10

			} else if board[row][0] == opponent {
				return -10
			}
		}
	}

	for col := 0; col < 3; col++ {
		if board[0][col] == board[1][col] &&
			board[1][col] == board[2][col] {
			if board[0][col] == player {
				return 10
			} else if board[0][col] == opponent {
				return -10
			}
		}
	}

	if board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		if board[0][0] == player {
			return 10
		} else if board[0][0] == opponent {
			return -10
		}
	}

	if board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		if board[0][2] == player {
			return +10
		} else if board[0][2] == opponent {
			return -10
		}
	}

	return 0
}

// EvaluateDynamic evaluates board of unlimited size from the view of player. -10 = loss / 10 = win / 0 = draw
func EvaluateDynamic(board [][]string, player string, opponent string) int {

	// Check rows
	for row := 0; row < len(board); row++ {
		if checkRow(board[row]) {
			if board[row][0] == player {
				return 10

			} else if board[row][0] == opponent {
				return -10
			}
		}
	}

	// Check columns
	for col := 0; col < len(board); col++ {
		previous := board[0][col]
		foundMatch := true
		for row := 0; row < len(board); row++ {
			if previous != board[row][col] {
				foundMatch = false
			}
			previous = board[row][col]

		}
		if foundMatch {
			if previous == player {
				return 10
			} else if previous == opponent {
				return -10
			}
		}
	}

	// Check cross
	previous := board[0][0]
	foundMatch := true
	for idx := 0; idx < len(board); idx++ {
		if previous != board[idx][idx] {
			foundMatch = false
		}
	}
	if foundMatch {
		if previous == player {
			return 10
		} else if previous == opponent {
			return -10
		}
	}

	previous = board[0][len(board)-1]
	foundMatch = true
	row := 0

	for col := len(board) - 1; col > -1; col-- {
		if previous != board[row][col] {
			foundMatch = false
		}
		previous = board[row][col]
		row++

	}
	if foundMatch {
		if previous == player {
			return 10
		} else if previous == opponent {
			return -10
		}
	}

	return 0
}

// Check if row has winning streak.
func checkRow(row []string) bool {
	previous := row[0]
	for col := 0; col < len(row); col++ {
		if previous != row[col] {
			return false
		}
		previous = row[col]
	}
	return true

}

// Find optimal move according to minimax
func minimax(board [][]string, depth int, isMax bool, player string, opponent string) int {
	score := Evaluate(board, player, opponent)

	if score == 10 {
		return score
	}

	if score == -10 {
		return score
	}

	if IsMovesLeft(board) == false {
		return 0
	}

	if isMax {
		best := -1000
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if board[i][j] == "-" {
					board[i][j] = player
					best = max(best, minimax(board, depth+1, !isMax, player, opponent))
					board[i][j] = "-"
				}
			}
		}
		return best
	}

	best := 1000

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "-" {
				board[i][j] = opponent
				best = min(best, minimax(board, depth+1, !isMax, player, opponent))
				board[i][j] = "-"
			}
		}
	}
	return best

}

//FindBestMove returns board with next best move for player (X/O).
func FindBestMove(strBoard string, player string, opponent string) string {
	bestVal := -1000
	bestMoveRow := -1
	bestMoveCol := -1
	board := ConvertStrToBoard(strBoard)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if board[i][j] == "-" {
				board[i][j] = player
				moveVal := minimax(board, 0, false, player, opponent)
				board[i][j] = "-"
				if moveVal > bestVal {
					bestMoveRow = i
					bestMoveCol = j
					bestVal = moveVal
				}
			}
		}
	}
	board[bestMoveRow][bestMoveCol] = player
	strBoard = ConvertBoardToStr(board)
	return strBoard
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// ConvertStrToBoard converts string presentation of board to that of two dimensional array.
func ConvertStrToBoard(strBoard string) [][]string {
	board := make([][]string, 3)
	board[0] = make([]string, 3)
	board[1] = make([]string, 3)
	board[2] = make([]string, 3)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			index := j + (i * 3)
			char := string(strBoard[index])
			board[i][j] = char

		}
	}
	return board
}

// ConvertBoardToStr converts two dimensional array presentation of board to that of string.
func ConvertBoardToStr(board [][]string) string {
	strBoard := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			char := board[i][j]
			strBoard += char
		}
	}
	return strBoard
}
