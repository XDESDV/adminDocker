package main

import (
	"adminDocker/app/server"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func newAdminDockerServer() error {

	// loading .env files in dev mode
	if os.Getenv("MODE") == "" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	srv := &server.AdminDocker{}

	srv.ParseParameters()

	// log format definition
	switch srv.LogFormat {
	case "HUMAN":
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "JSON":
		// Already default
	default:
		log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
	}

	// setup router
	srv.Router = setupRouter()

	// api.SetupRouter(srv.Router)

	server.SetServer(srv)

	return nil
}
