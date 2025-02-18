package common

import (
	"adminDocker/app/models"
	"adminDocker/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	messageTypes := &models.MessageTypes{
		OK: "ping.Done",
	}
	SendResponse(c, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "pong"))
}

func Version(c *gin.Context) {
	srv := server.GetServer()
	messageTypes := &models.MessageTypes{
		OK: "version.Done",
	}
	SendResponse(c, http.StatusOK, models.Success(http.StatusOK, messageTypes.OK, "Version:"+srv.Version))
}

func SendResponse(c *gin.Context, status int, response interface{}) {
	c.JSON(status, response)
}
