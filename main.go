package main

import (
	"log"
	"net/http"
	"os"

	"github.com/artoju/tic-tac-toe/config"
	"github.com/artoju/tic-tac-toe/db"
	"github.com/artoju/tic-tac-toe/game/state"
	"github.com/artoju/tic-tac-toe/handlers"
	"github.com/artoju/tic-tac-toe/redis"
	"github.com/artoju/tic-tac-toe/socket"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "DELETE", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		Debug:            false,
	})
	handler := c.Handler(router)

	conf, err := config.Init()
	if err != nil {
		panic(err)
	}
	sh := handlers.StateHandler{}

	srv := &http.Server{
		Addr:    ":" + conf.Server.Port,
		Handler: handler,
	}

	if conf.Server.StateHandler == "db" {
		ddb, err := db.Init(conf)
		if err != nil {
			panic(err)
		}
		dbh := state.DBStateHandler{DB: ddb}
		sh.GameState = dbh
	} else if conf.Server.StateHandler == "redis" {
		client, err := redis.Init(conf)
		if err != nil {
			panic(err)
		}
		dbh := state.RedisStateHandler{DB: client}
		sh.GameState = dbh
	} else {
		f := state.FileStateHandler{Filepath: conf.FileHandler.Path}
		sh.GameState = f
		if _, err := os.Stat(conf.FileHandler.Path); os.IsNotExist(err) {
			os.Mkdir(conf.FileHandler.Path, 0777)
		}
	}
	sh.IsSinglePlayer = false

	router.HandleFunc("/api/v1/games", sh.GetGamesHandler).Methods("GET")
	router.HandleFunc("/api/v1/games", sh.StartGameHandler).Methods("POST")

	router.HandleFunc("/api/v1/games/{id}", sh.GetGameHandler).Methods("GET")
	router.HandleFunc("/api/v1/games/{id}", sh.UpdateGameHandler).Methods("PUT")
	router.HandleFunc("/api/v1/games/{id}", sh.DeleteGameHandler).Methods("DELETE")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("SERVER ONLINE"))
	})

	socket.MainLobby.GameStateHandler = &sh.GameState
	go socket.MainLobby.Run()

	router.HandleFunc("/ws", socket.LobbyHandler)
	router.HandleFunc("/ws/{id}/{token}", socket.GameHandler)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
