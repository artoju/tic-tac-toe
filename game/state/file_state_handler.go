package state

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type FileStateHandler struct {
	Filepath string
}

func (f FileStateHandler) CreateState(req Game, c chan GameResponse) {
	gameFile := f.Filepath + req.ID + ".txt"

	info, err := os.Stat(gameFile)
	if !os.IsNotExist(err) {
		c <- GameResponse{
			Error: errors.New("Game already exists"),
		}
		return
	}
	if info.IsDir() {
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	file, err := os.OpenFile(gameFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		log.Println("Error opening file: " + err.Error())
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}

	str := req.ID + ":" + req.Board + ":" + req.Status + ":" + req.NextPlayer + "\n"
	_, err = file.WriteString(str)
	if err != nil {
		log.Println("Error writing file: " + err.Error())
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
func (f FileStateHandler) GetState(req Game, c chan GameResponse) {

	gameFile := f.Filepath + req.ID + ".txt"

	info, err := os.Stat(gameFile)
	if os.IsNotExist(err) {
		c <- GameResponse{
			Error: errors.New("No game found for ID " + gameFile),
		}
		return
	}
	if info.IsDir() {
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}

	file, err := os.OpenFile(gameFile, os.O_APPEND|os.O_RDONLY, 0600)
	defer file.Close()
	if err != nil {
		log.Println("Error opening file: " + err.Error())
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}

	scanner := bufio.NewScanner(file)

	line := ""
	for scanner.Scan() {
		line = scanner.Text()
	}
	game := Game{}
	if line == "" {
		game.ID = req.ID
		game.Board = "---------"
		game.Status = "RUNNING"
	} else {
		parts := strings.Split(line, ":")
		game.ID, game.Board, game.Status, game.NextPlayer = parts[0], parts[1], parts[2], parts[3]
	}

	c <- GameResponse{
		Error: nil,
		Game:  game,
	}

}

func (f FileStateHandler) SaveState(req Game, c chan GameResponse) {
	gameFile := f.Filepath + req.ID + ".txt"
	file, err := os.OpenFile(gameFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer file.Close()
	if err != nil {
		log.Println("Error opening file for update: " + err.Error())
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}

	str := req.ID + ":" + req.Board + ":" + req.Status + ":" + req.NextPlayer + "\n"
	_, err = file.WriteString(str)
	if err != nil {
		log.Println("Error writing file: " + err.Error())
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	game := Game{
		ID:         req.ID,
		Board:      req.Board,
		Status:     req.Status,
		NextPlayer: req.NextPlayer,
	}
	c <- GameResponse{
		Game:  game,
		Error: nil,
	}

}

func (f FileStateHandler) DeleteState(id string, c chan GameResponse) {

	gameFile := f.Filepath + id + ".txt"

	info, err := os.Stat(gameFile)
	if os.IsNotExist(err) {
		c <- GameResponse{
			Error: errors.New("No game found for ID"),
		}
		return
	}
	if info.IsDir() {
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	err = os.Remove(gameFile)
	if err != nil {
		log.Println("Error getting working directory: " + err.Error())
		c <- GameResponse{
			Error: errors.New("Internal error"),
		}
		return
	}
	c <- GameResponse{
		Error: nil,
	}

}

func (f FileStateHandler) GetStates() ([]Game, error) {
	games := make([]Game, 0)

	gameDir := f.Filepath

	files, err := ioutil.ReadDir(gameDir)
	if err != nil {
		log.Println("Error reading directory: " + err.Error())
		return games, errors.New("Internal error")
	}

	for _, f := range files {
		gameFile := gameDir + f.Name()

		file, err := os.OpenFile(gameFile, os.O_APPEND|os.O_RDONLY, 0600)
		defer file.Close()
		if err != nil {
			log.Println("Error opening file: " + err.Error())
			return games, errors.New("Internal error")
		}

		scanner := bufio.NewScanner(file)

		line := ""
		for scanner.Scan() {
			line = scanner.Text()
		}
		game := Game{}
		parts := strings.Split(line, ":")
		game.ID, game.Board, game.Status = parts[0], parts[1], parts[2]

		games = append(games, game)
	}
	return games, nil
}

func (f FileStateHandler) ValidatePlayer(token string, id string) (*string, error) {
	return nil, nil
}
