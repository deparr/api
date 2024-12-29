package main

import (
	"log/slog"

	"os"

	"github.com/deparr/api/pkg/server"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("loading environ:", "error", err)
		os.Exit(1)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	err = server.ListenAndServe(host, port)
	if err != nil {
		slog.Error("failed to start server:", "error", err)
		os.Exit(1)
	}
}
