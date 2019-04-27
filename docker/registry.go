package docker

import (
	"fmt"
	"strings"

	"github.com/heroku/docker-registry-client/registry"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/utils"
)

// Check if image (project name) by user (username) with tag exists in
// the private repository. The hash must be
func (c *Client) hubHasImage(username, name, hash string) (bool, error) {
	if len(hash) < 8 {
		return false, nil
	}

	repo := fmt.Sprintf("%s/%s", username, name)
	repos, err := c.hub.Repositories()

	if err != nil {
		return false, errors.Wrap(err, "error getting repos from private registry")
	}
	if !utils.StrIsIn(repo, repos) {
		// repos doesn't even exist. Image does not exist!
		return false, nil
	}

	tags, err := c.hub.Tags(repo)
	if err != nil {
		return false, errors.Wrapf(err, "error getting tags from private registry (%s) for repo (%s)", c.hub.URL, repo)
	}
	hash = strings.ToLower(hash)
	for _, t := range tags {
		if strings.HasPrefix(strings.ToLower(t), hash) {
			return true, nil
		}
	}
	return false, nil
}

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
func formRegistryTag(username, name, hash string) string {
	addr := viper.GetString("registry.domain")
	port := viper.GetInt("registry.port")
	if port != 0 && port != 80 && port != 443 {
		addr = fmt.Sprintf("%s:%d", addr, port)
	}

	return strings.ToLower(fmt.Sprintf("%s/%s/%s:%s", addr, username, name, hash))
}
