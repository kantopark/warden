package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/utils"
)

// A cli for the docker images. Each function and it's various RunTags are
// essentially a Docker Image which the cluster will deploy. The cli ensures
// that the required image is in the repository so that the various nodes are
// able to get the images as required.
type Client struct {
	cli   *client.Client
	ctx   context.Context
	hub   *registry.Registry
	redis *redis.Client
}

const (
	redisContainer = "warden_redis"
)

// Creates a new cli to oversee operations of the Docker cli
func NewClient() (*Client, error) {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new Docker Client")
	}

	hub, err := newHub()
	if err != nil {
		return nil, err
	}

	// Login to docker hub if username is provided. This is probably not needed unless the default registry
	// is credential secured
	if !utils.StrIsEmptyOrWhitespace(viper.GetString("docker.username")) {
		_, err := c.RegistryLogin(ctx, types.AuthConfig{
			Username:      viper.GetString("docker.username"),
			Password:      viper.GetString("docker.password"),
			Email:         viper.GetString("docker.email"),
			ServerAddress: viper.GetString("docker.serveraddr"),
		})
		if err != nil {
			return nil, errors.Wrap(err, "error logging in")
		}
	}

	cli := &Client{
		cli: c,
		ctx: ctx,
		hub: hub,
	}
	if err := cli.startRedis(); err != nil {
		return nil, err
	}

	return cli, nil
}

func (c *Client) startRedis() error {
	redisImage := viper.GetString("redis.image")
	redisHost := viper.GetString("redis.addr")
	redisPort := viper.GetString("redis.port")

	// pull the redis image if we can't find it in the local repo
	if img, _ := c.FindImageByName(redisImage); img == nil {
		err := c.PullImage(redisImage, &ImagePullOptions{UseDockerHub: true})
		if err != nil {
			return errors.Wrap(err, "error pulling Redis image")
		}
	}

	if viper.GetBool("redis.restart_if_exist") {
		if err := c.removeRedis(); err != nil {
			return err
		}
	}
	if err := c.runRedisContainer(); err != nil {
		return err
	}

	// start up redis client
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	c.redis = redis.NewClient(&redis.Options{
		Addr:        redisAddr,
		Password:    viper.GetString("redis.password"),
		DB:          viper.GetInt("redis.DB"),
		IdleTimeout: 5 * time.Minute,
	})

	if _, err := c.redis.Ping().Result(); err != nil {
		return errors.Wrapf(err, "error connecting to the redis server at %s", redisAddr)
	}

	return nil
}

// Teardowns the Client object properly
func (c *Client) Close() (err error) {
	if viper.GetBool("redis.remove_on_exit") {
		err = c.removeRedis()
	} else {
		err = c.redis.FlushAll().Err()
	}
	if err != nil {
		return err
	}

	if err := c.cli.Close(); err != nil {
		return errors.Wrap(err, "error stopping docker cli")
	}
	c.ctx.Done()

	// Remove image build folder. This folder houses the tar folders that contain the files
	// to build the image
	addr := viper.GetString("registry.domain")
	port := viper.GetInt("registry.port")
	if port != 0 {
		addr = fmt.Sprintf("%s-%d", addr, port)
	}
	if utils.PathExists(addr) {
		os.RemoveAll(addr)
	}

	return nil
}

// Kills and remove the redis container if it exists
func (c *Client) removeRedis() error {
	// search for containers with the defined redis container name
	ftr := filters.NewArgs()
	ftr.Add("name", redisContainer)
	containers, _ := c.cli.ContainerList(
		c.ctx,
		types.ContainerListOptions{
			Filters: ftr,
			All:     true})

	// if container exist (there should only be 1), remove it
	if len(redisContainer) > 0 {
		for _, con := range containers {
			if err := c.cli.ContainerRemove(c.ctx, con.ID, types.ContainerRemoveOptions{
				Force: true,
			}); err != nil {
				return errors.Wrap(err, "error removing Redis container")
			}
		}
	}
	return nil
}

// Starts a new redis container
func (c *Client) runRedisContainer() error {
	redisImage := viper.GetString("redis.image")
	redisHost := viper.GetString("redis.addr")
	redisPort := viper.GetString("redis.port")

	ftr := filters.NewArgs()
	ftr.Add("name", redisContainer)
	containers, _ := c.cli.ContainerList(
		c.ctx,
		types.ContainerListOptions{
			Filters: ftr,
			All:     true})

	if len(containers) > 0 {
		log.Println("Instance of redis already running")
		return nil
	}

	log.Println("starting new redis container")
	redisCon, err := c.cli.ContainerCreate(
		c.ctx,
		&container.Config{
			Image:        redisImage,
			ExposedPorts: nat.PortSet{nat.Port(redisPort): struct{}{}},
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(redisPort): {{HostIP: redisHost, HostPort: redisPort}}},
		},
		nil,
		redisContainer,
	)
	if err != nil {
		return errors.Wrap(err, "error creating Redis container")
	}

	if err = c.cli.ContainerStart(c.ctx, redisCon.ID, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "error starting Redis container")
	}

	return nil
}
