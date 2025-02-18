package main

import (
	controllers "adminDocker/app/controllers/common"
	routes "adminDocker/app/routes/common"

	"github.com/gin-gonic/gin"
)

// init the router
func setupRouter() *gin.Engine {
	router := routes.SetupRouter()
	router.GET("/ping", controllers.Ping)
	router.GET("/version", controllers.Version)

	return router
}
