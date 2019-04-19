package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

// A cli for the docker images. Each function and it's various RunTags are
// essentially a Docker Image which the cluster will deploy. The cli ensures
// that the required image is in the repository so that the various nodes are
// able to get the images as required.
type Client struct {
	cli     *client.Client
	ctx     context.Context
	redisId string
}

// Creates a new cli to oversee operations of the Docker cli
func NewClient() (*Client, error) {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new Docker Client")
	}

	cli := &Client{
		cli: c,
		ctx: ctx,
	}
	if err := cli.startRedis(); err != nil {
		return nil, err
	}

	return cli, nil
}

func (c *Client) startRedis() error {
	redisImage := "redis:latest"
	containerName := "warden-redis"

	imagePullResp, err := c.cli.ImagePull(
		c.ctx,
		fmt.Sprintf("docker.io/library/%s", redisImage),
		types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error pulling Redis image")
	}
	defer imagePullResp.Close()
	streamResponse(imagePullResp)

	// Remove redis containers that were somehow not destroyed previously
	ftr := filters.NewArgs()
	ftr.Add("name", containerName)
	redisContainer, _ := c.cli.ContainerList(
		c.ctx,
		types.ContainerListOptions{
			Filters: ftr,
			All:     true})

	if len(redisContainer) > 0 {
		for _, _container := range redisContainer {
			if err := c.removeRedis(_container.ID); err != nil {
				log.Fatalln(err)
			}
		}
	}

	// Starts a new redis container
	redisCont, err := c.cli.ContainerCreate(
		c.ctx,
		&container.Config{
			Image:        redisImage,
			ExposedPorts: nat.PortSet{"6379": struct{}{}},
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{nat.Port("6379"): {{HostIP: "127.0.0.1", HostPort: "6379"}}},
		},
		nil,
		"warden-redis",
	)
	if err != nil {
		return errors.Wrap(err, "error creating Redis container")
	}

	c.redisId = redisCont.ID
	err = c.cli.ContainerStart(c.ctx, c.redisId, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "error starting Redis container")
	}

	return nil
}

// Teardowns the Client object properly
func (c *Client) Close() error {
	if err := c.removeRedis(c.redisId); err != nil {
		return err
	}
	if err := c.cli.Close(); err != nil {
		return errors.Wrap(err, "error stopping docker cli")
	}
	c.ctx.Done()
	return nil
}

// Kills and remove the redis container
func (c *Client) removeRedis(containerId string) error {
	if err := c.cli.ContainerRemove(c.ctx, containerId, types.ContainerRemoveOptions{
		Force: true,
	}); err != nil {
		return errors.Wrap(err, "error removing Redis container")
	}
	return nil
}
