package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	CHUNK_SIZE     int
	BROADCAST_PORT int
}

var Config AppConfig

func NewConfig(CHUNK_SIZE int, BROADCAST_PORT int) *AppConfig {
	return &AppConfig{
		CHUNK_SIZE:     CHUNK_SIZE,
		BROADCAST_PORT: BROADCAST_PORT,
	}
}

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chunkSize, err := strconv.ParseInt(os.Getenv("CHUNK_SIZE"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing CHUNK_SIZE:", err)
	}

	broadcastPort, err := strconv.ParseInt(os.Getenv("BROADCAST_PORT"), 10, 64)
	if err != nil {
		log.Fatal("Error parsing BROADCAST_PORT:", err)
	}

	Config = AppConfig{
		CHUNK_SIZE:     int(chunkSize),
		BROADCAST_PORT: int(broadcastPort),
	}
}
