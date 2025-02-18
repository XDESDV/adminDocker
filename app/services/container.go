package services

import (
	"adminDocker/app/server"
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type Container struct {
	clientDocker *client.Client
	validate     *validator.Validate
	logs         *zerolog.Logger
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

func (c *Container) Close() {
	c.clientDocker.Close()
}
