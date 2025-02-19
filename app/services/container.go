package services

import (
	"adminDocker/app/server"
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"io"
)

type Container struct {
	clientDocker *client.Client
	validate     *validator.Validate
	logs         *zerolog.Logger
}

// Structure the stats CPU et RAM of container
type ContainerStats struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   string  `json:"memory_usage"`
	MemoryLimit   string  `json:"memory_limit"`
	MemoryPercent float64 `json:"memory_percent"`
}

func NewServiceContainer(logs *zerolog.Logger) (*Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Container{
		clientDocker: cli,
		validate:     validator.New(),
		logs:         logs,
	}, nil
}

func (c *Container) ListDocker() ([]types.Container, error) {
	if server.GetServer().DockerFake || c.clientDocker == nil {
		// Retourner des données factices si Docker n'est pas disponible
		fakeContainers := []types.Container{
			{
				ID:     "123456789abc",
				Names:  []string{"/fake-nginx"},
				State:  "running",
				Status: "Up 10 minutes",
			},
			{
				ID:     "987654321xyz",
				Names:  []string{"/fake-redis"},
				State:  "exited",
				Status: "Exited (0) 5 minutes ago",
			},
		}
		c.logs.Warn().Msg("Retour de données fake.")
		return fakeContainers, nil
	}

	containers, err := c.clientDocker.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		c.logs.Error().Err(err).Msg("")
		return nil, err
	}
	return containers, nil

}

// Stop a contanier
func (c *Container) StopDocker(containerID string) (bool, error) {
	if server.GetServer().DockerFake || c.clientDocker == nil {
		return true, nil
	}

	err := c.clientDocker.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		c.logs.Error().Err(err).Msg("")
		return false, err
	}

	return true, nil
}

// Start a contanier
func (c *Container) StartDocker(containerID string) (bool, error) {
	if server.GetServer().DockerFake || c.clientDocker == nil {
		return true, nil
	}

	err := c.clientDocker.ContainerStart(context.Background(), containerID, container.StartOptions{})
	if err != nil {
		c.logs.Error().Err(err).Msg("")
		return false, err
	}

	return true, nil
}

// Pull image
func (c *Container) PullImage(imageName string) error {
	_, _, err := c.clientDocker.ImageInspectWithRaw(context.Background(), imageName)
	if err == nil {
		c.logs.Info().Msgf("Image %s found locally", imageName)
		return nil // Image already exists
	}

	c.logs.Info().Msgf("Pulling image: %s", imageName)
	out, err := c.clientDocker.ImagePull(context.Background(), imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer out.Close()
	io.Copy(io.Discard, out) // Consume output to avoid blocking
	return nil
}

// Create a container
func (c *Container) CreateDocker(containerName, imageName string, command []string, ports []string) (string, error) {
	// Configure port bindings
	portBindings := map[nat.Port][]nat.PortBinding{}
	for _, port := range ports {
		portBindings[nat.Port(port+"/tcp")] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: port},
		}
	}

	resp, err := c.clientDocker.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: imageName,
			Cmd:   command,
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil, nil, containerName,
	)

	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return resp.ID, nil
}

// Stream container's logs
func (c *Container) StreamContainerLogs(containerID string, output io.Writer) error {
	logs, err := c.clientDocker.ContainerLogs(context.Background(), containerID, container.LogsOptions{})
	if err != nil {
		return fmt.Errorf("failed to retrieve logs: %w", err)
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		fmt.Fprintln(output, scanner.Text())
	}

	return scanner.Err()
}

// GetContainerStats gets the CPU/RAM stats  of a container
func (c *Container) GetContainerStats(containerID string) (*ContainerStats, error) {
	stats, err := c.clientDocker.ContainerStats(context.Background(), containerID, false)
	if err != nil {
		return nil, err
	}
	defer stats.Body.Close()

	var statsJSON types.StatsJSON
	body, err := io.ReadAll(stats.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &statsJSON); err != nil {
		return nil, err
	}

	// Calculate CPU
	cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - statsJSON.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
	var cpuPercent float64 = 0
	if systemDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) * 100.0
	}

	// Calculate RAM
	memoryUsage := float64(statsJSON.MemoryStats.Usage)
	memoryLimit := float64(statsJSON.MemoryStats.Limit)

	// Format the values
	statsResult := &ContainerStats{
		CPUPercent:  cpuPercent,
		MemoryUsage: fmt.Sprintf("%.2f MiB", memoryUsage/1024/1024),
		MemoryLimit: fmt.Sprintf("%.2f MiB", memoryLimit/1024/1024),
	}

	return statsResult, nil
}

// Get container's name
func (c *Container) GetContainerName(containerID string) (*types.ContainerJSON, error) {
	containerJSON, err := c.clientDocker.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, err
	}
	return &containerJSON, nil
}
func (c *Container) Close() {
	c.clientDocker.Close()
}
