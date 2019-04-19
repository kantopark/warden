package docker

import (
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// A manager for the docker images. Each function and it's various RunTags are
// essentially a Docker Image which the cluster will deploy. This manager ensures
// that the required image is in the repository so that the various nodes are
// able to get the images as required.
type Manager struct {
	client *client.Client
}

func NewManager() (*Manager, error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new Docker Client")
	}
	return &Manager{
		client: c,
	}, nil
}
