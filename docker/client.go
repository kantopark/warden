package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

// A cli for the docker images. Each function and it's various RunTags are
// essentially a Docker Image which the cluster will deploy. The cli ensures
// that the required image is in the repository so that the various nodes are
// able to get the images as required.
type Client struct {
	cli *client.Client
	ctx context.Context
}

var redisId string

// Creates a new cli to oversee operations of the Docker cli
func NewClient() (*Client, error) {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new Docker Client")
	}

	if err := startRedis(c, ctx); err != nil {
		return nil, err
	}

	return &Client{
		cli: c,
		ctx: ctx,
	}, nil
}

func startRedis(c *client.Client, ctx context.Context) error {
	redisImage := "redis:latest"

	imagePullResp, err := c.ImagePull(ctx, redisImage, types.ImagePullOptions{})
	if err != nil {
		return errors.Wrap(err, "error pulling Redis image")
	}
	defer imagePullResp.Close()
	streamResponse(imagePullResp)

	redisCont, err := c.ContainerCreate(
		ctx,
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

	redisId = redisCont.ID
	err = c.ContainerStart(ctx, redisId, types.ContainerStartOptions{})
	if err != nil {
		return errors.Wrap(err, "error starting Redis container")
	}

	return nil
}

// Teardowns the Client object properly
func (c *Client) Close() error {
	// Stops and remove redis container
	if err := c.cli.ContainerKill(c.ctx, redisId, "KILL"); err != nil {
		return errors.Wrap(err, "error killing Redis container")
	}
	if err := c.cli.ContainerRemove(c.ctx, redisId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	}); err != nil {
		return errors.Wrap(err, "error removing Redis container")
	}
	if err := c.cli.Close(); err != nil {
		return errors.Wrap(err, "error stopping docker cli")
	}
	c.ctx.Done()
	return nil
}
