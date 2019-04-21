package docker

import (
	"fmt"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Creates a new registry hub
func newHub() (*registry.Registry, error) {
	// Creates a connection to local registry
	regAddr := fmt.Sprintf("%s://%s", viper.GetString("registry.protocol"), viper.GetString("registry.domain"))
	regPort := viper.GetInt("registry.port")
	if regPort != 0 && regPort != 80 && regPort != 443 {
		regAddr = fmt.Sprintf("%s:%d", regAddr, regPort)
	}
	hub, err := registry.New(regAddr, viper.GetString("registry.username"), viper.GetString("registry.password"))
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to registry")
	}

	return hub, nil
}

// Creates the tag for the image. Name should refer to the project name and alias
// refer to the alias of the built image.
func formRegistryTag(username, name, tag string) string {
	addr := viper.GetString("registry.domain")
	port := viper.GetInt("registry.port")
	if port != 0 && port != 80 && port != 443 {
		addr = fmt.Sprintf("%s:%d", addr, port)
	}

	return strings.ToLower(fmt.Sprintf("%s/%s/%s:%s", addr, username, name, tag))
}
