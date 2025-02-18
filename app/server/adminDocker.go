package server

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var server *AdminDocker

// AdminDocker Structure
type AdminDocker struct {
	Router    *gin.Engine
	Version   string
	Port      string
	TokenKey  string
	Origin    string
	LogFormat string
	Mode      string
}

func (a *AdminDocker) ParseParameters() {
	a.LogFormat = os.Getenv("LOG_FORMAT")
	a.Version = os.Getenv("API_VERSION")
	a.Port = os.Getenv("API_PORT")
	a.TokenKey = os.Getenv("TOKEN_KEY")
	a.Origin = os.Getenv("ALLOW_ORIGIN")
	a.Mode = os.Getenv("MODE")
}

// ListenAndServe listens on the TCP network address addr and then calls Serve with handler to handle requests on incoming connections.
// https://github.com/gin-gonic/gin
func (a *AdminDocker) ListenAndServe() error {
	srv := &http.Server{
		Addr:              a.Port,
		Handler:           a.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// start
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("Unable to listen and serve: %v", err)
		return err
	}
	return nil
}

// SetServer init mongo database
func SetServer(s *AdminDocker) {
	server = s
}

// GetServer Flashcards
func GetServer() *AdminDocker {
	return server
}
