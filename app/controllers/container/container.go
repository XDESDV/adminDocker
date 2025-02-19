package container

import (
	"adminDocker/app/controllers/common"
	"adminDocker/app/models"
	"adminDocker/app/services"
	"errors"
	"fmt"
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

// Stop controller to stop a container
func (c *Container) Stop(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)

	messageTypes := &models.MessageTypes{
		OK:                  "container.Start.Found",
		InternalServerError: "container.Start.Error",
	}

	containerID := ctx.Param("id")

	Ok, err := c.containerService.StopDocker(containerID)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
	}
	if !Ok {
		common.SendResponse(ctx, http.StatusInternalServerError, "The container wasn't stopped")
	}
	common.SendResponse(ctx, http.StatusOK, "The container was stopped")
}

// Start controller to start a container
func (c *Container) Start(ctx *gin.Context) {
	var params models.QueryParams

	params.Parse(ctx)

	messageTypes := &models.MessageTypes{
		OK:                  "container.Start.Found",
		InternalServerError: "container.Start.Error",
	}

	containerID := ctx.Param("id")

	Ok, err := c.containerService.StartDocker(containerID)
	if err != nil {
		common.SendResponse(ctx, http.StatusInternalServerError, models.KnownError(http.StatusInternalServerError, messageTypes.InternalServerError, err))
	}
	if !Ok {
		common.SendResponse(ctx, http.StatusInternalServerError, "The container wasn't started")
	}
	common.SendResponse(ctx, http.StatusOK, "The container was started")
}

// Create controller to create a container
func (c *Container) Create(ctx *gin.Context) {
	// Get the name of the image from param
	containerName := ctx.Param("name")

	// Structure the request
	var req struct {
		Image   string   `json:"image" binding:"required"`
		Command []string `json:"command"`
		Ports   []string `json:"ports"`
	}

	// Handle request error
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Pull Image if not exists
	if err := c.containerService.PullImage(req.Image); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create container
	containerID, err := c.containerService.CreateDocker(containerName, req.Image, req.Command, req.Ports)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Start container
	ok, err := c.containerService.StartDocker(containerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start container"})
		return
	}

	// Send response headers
	ctx.Writer.Header().Set("Content-Type", "text/plain")
	ctx.Writer.WriteHeader(http.StatusOK)

	// Display container information
	fmt.Fprintf(ctx.Writer, "Starting Nginx container...\n")
	fmt.Fprintf(ctx.Writer, "Container ID: %s\n", containerID)
	fmt.Fprintf(ctx.Writer, "Logs:\n")

	// Stream logs in real-time
	c.containerService.StreamContainerLogs(containerID, ctx.Writer)

}

// Get Container Resources
func (c *Container) GetContainerResources(ctx *gin.Context) {
	containerID := ctx.Param("id")

	stats, err := c.containerService.GetContainerStats(containerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the name of the container
	containerInfo, err := c.containerService.GetContainerName(containerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Formater la sortie
	response := fmt.Sprintf(
		"Conteneur : %s\nCPU : %.2f%%\nMémoire : %s / %s",
		containerInfo.Name, stats.CPUPercent, stats.MemoryUsage, stats.MemoryLimit,
	)

	ctx.String(http.StatusOK, response)
}
