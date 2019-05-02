package docker

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"

	"warden/config"
)

var cli *Client

func init() {
	config.ReadConfig()
	_cli, err := NewClient()
	if err != nil {
		log.Fatalln(err)
	}
	cli = _cli
}

func TestNewClient(t *testing.T) {
	_, err := NewClient()
	assert.Nil(t, err)
}
