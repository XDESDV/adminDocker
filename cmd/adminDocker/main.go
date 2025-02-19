package main

import (
	"adminDocker/app/routes/dockers"
	"adminDocker/app/server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"os"

	"github.com/rs/zerolog/log"

	controllers "adminDocker/app/controllers/common"
	routes "adminDocker/app/routes/common"
)

func main() {
	if err := newAdminDockerServer(); err != nil {
		log.Fatal().Err(err).Msg("Unable to create new server")
		os.Exit(51)
	}
	log.Debug().Msg("API launched with human readable log")

	srv := server.GetServer()
	srv.ListenAndServe()
}

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

	err := dockers.SetupRouter(srv.Router, &log.Logger)
	if err != nil {
		return err
	}
	server.SetServer(srv)

	return nil
}

func setupRouter() *gin.Engine {
	router := routes.SetupRouter()
	router.GET("/ping", controllers.Ping)
	router.GET("/version", controllers.Version)

	return router
}
