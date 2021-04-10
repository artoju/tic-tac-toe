package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// TicTacToe game server configurations.
// Only one handler should be initialized
// for game state handling.
type Config struct {
	Server          Server          `yaml:"server"`
	DatabaseHandler DatabaseHandler `yaml:"databaseHandler"`
	FileHandler     FileHandler     `yaml:"fileHandler"`
	RedisHandler    RedisHandler    `yaml:"redisHandler"`
}

type Server struct {
	Port string `yaml:"port"`
	// Selected game state handler. Options are: redis / file / db
	StateHandler string `yaml:"stateHandler"`
}

// DatabaseHandler contains required credentials
// for AWS DynamoDB connection.
type DatabaseHandler struct {
	KeyId     string `yaml:"keyId"`
	SecretKey string `yaml:"secretKey"`
	TableName string `yaml:"tableName"`
	Region    string `yaml:"region"`
}

// FileHandler contains path for saving
// game state to file.
type FileHandler struct {
	Path string `yaml:"path"`
}

// RedisHandler contains required credentials
// for Redis connection.
type RedisHandler struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
}

func Init() (*Config, error) {
	configFile, err := ioutil.ReadFile("config/config.yml")
	if err != nil {
		panic(err)
	}
	var conf Config
	err = yaml.UnmarshalStrict(configFile, &conf)
	if err != nil {
		panic(err)
	}
	return &conf, nil
}
