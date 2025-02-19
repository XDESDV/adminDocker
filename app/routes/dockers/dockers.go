package dockers

import (
	controller "adminDocker/app/controllers/container"
	services "adminDocker/app/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func SetupRouter(g *gin.Engine, logs *zerolog.Logger) error {

	containerService, err := services.NewServiceContainer(logs)
	if err != nil {
		return err
	}
	containerController := controller.New(containerService, logs)

	v1 := g.Group("/v1")
	{
		dockersV1 := v1.Group("/dockers")
		{
			dockersV1.GET("", containerController.Get)
			dockersV1.POST("/stop/:id", containerController.Stop)
			dockersV1.POST("/start/:id", containerController.Start)
			dockersV1.PUT("/create/:name", containerController.Create)
			dockersV1.GET("/:id/ressources", containerController.GetContainerResources)
		}
	}

	return nil
}
