package common

import (
	"adminDocker/app/controllers/common"
	"adminDocker/app/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// InitialiseRouter initialization of web service routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
	noRoute(router)
	useCORS(router)
	return router
}

func useCORS(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		allowOrigin := os.Getenv("ALLOW_ORIGIN")
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Manage OPTIONS queries, used for CORS preflighting
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	})
}

func noRoute(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		messageTypes := &models.MessageTypes{
			NotFound: "Ressource.NotFound",
		}
		common.SendResponse(c, http.StatusNotFound, models.Success(http.StatusNotFound, messageTypes.NotFound, "ressource not found"))

	})
}
