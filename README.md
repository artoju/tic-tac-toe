# tic-tac-toe

Trying out websockets and redis as well as other game state saving methods on Go. Play tic-tac-toe against CPU or other players.
## Running project

Get dependencies
```
go get -d -v
```
Build executable
```
go build main.go
```
Copy and or edit configurations
```
cp config/example.config.yml config/config.yml
```
Run app
```
./main
```
Or run the app with docker using my image from dockerhub
```
docker-compose up -d