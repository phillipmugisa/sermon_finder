package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/phillipmugisa/sermon_finder/server"
)

func main() {
	logger := slog.New(slog.Default().Handler())

	envErr := godotenv.Load()
	if envErr != nil {
		logger.Error("no env file found")
	}
	PORT := os.Getenv("PORT")
	if PORT == "" {
		logger.Error("no port declaration found in env file")
	}

	server := server.NewServer(PORT, nil, logger)
	err := server.Start()
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}
