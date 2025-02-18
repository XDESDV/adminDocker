package container

import (
	"adminDocker/app/controllers/common"
	"adminDocker/app/models"
	"adminDocker/app/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Container struct {
	containerService *services.Container
	logs             *zerolog.Logger
}

func New(containerService *services.Container, logs *zerolog.Logger) *Container {
	return &Container{
		containerService: containerService,
		logs:             logs,
	}
}

// Get controller to get list of containers
func (c *Container) Get(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)
	messageTypes := &models.MessageTypes{
		OK:                  "container.Search.Found",
		BadRequest:          "container.Search.BadRequest",
		NotFound:            "container.Search.NotFound",
		InternalServerError: "container.Search.Error",
	}

	containers, err := c.containerService.ListDocker()
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
	}
	totalCount := len(containers)
	if totalCount == 0 {
		status := http.StatusNotFound
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Data not found. ")))
	}

	low := params.Offset - 1
	if low == -1 {
		low = 0
	}

	// Available CountMax calculation
	maxCount := params.Count
	if maxCount == 0 {
		maxCount = 100
	}

	high := maxCount + low
	if high > totalCount {
		high = totalCount
	}

	if low > high {
		status := http.StatusBadRequest
		common.SendResponse(ctx, status, models.KnownError(status, messageTypes.NotFound, errors.New(" Offset cannot be higher than count. ")))
	}

	sendingContainers := containers[low:high]

	meta := models.MetaResponse{
		ObjectName: "Dockers",
		TotalCount: totalCount,
		Count:      len(sendingContainers),
		Offset:     low + 1,
	}

	response := &models.WSResponse{
		Meta: meta,
		Data: sendingContainers,
	}

	common.SendResponse(ctx, http.StatusOK, response)
}
