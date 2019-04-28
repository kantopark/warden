package deploy

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/utils"
)

type Manager interface {
	Close() error
	DeployInstance(d Deployment) error
	StopInstance(d Deployment) error
}

func NewManager() (Manager, error) {
	r := strings.TrimSpace(strings.ToLower(viper.GetString("deploy.type")))
	switch r {
	case "docker":
		return newDockerRunner()
	default:
		return nil, errors.Errorf("Unknown deploy type: %s", r)
	}
}

type Deployment struct {
	Project    string
	Hash       string
	MinReplica int
	MaxReplica int
}

// Validate and set sane defaults for the Deployment object
func (d *Deployment) validate() error {
	var err []string

	if utils.StrIsEmptyOrWhitespace(d.Project) {
		err = append(err, "Project field in deployment must be specified")
	}
	if utils.StrIsEmptyOrWhitespace(d.Hash) {
		err = append(err, "Hash field in deployment must be specified")
	}

	// setting up min and max replicas
	if d.MinReplica == 0 && d.MaxReplica == 0 {
		d.MaxReplica = 1
		d.MinReplica = 1
	}

	if d.MaxReplica < 0 || d.MinReplica < 0 {
		err = append(err, "Min/Max Replicas must be >= 0")
	}

	if d.MaxReplica < d.MinReplica {
		err = append(err, "Max replica must be >= min replica")
	}

	// returns errors if errors exist
	if len(err) > 0 {
		return errors.New(strings.Join(err, "\n"))
	}
	return nil
}

func (d *Deployment) ImageName() string {
	addr := viper.GetString("registry.domain")
	port := viper.GetInt("registry.port")
	if port != 0 && port != 80 && port != 443 {
		addr = fmt.Sprintf("%s:%d", addr, port)
	}

	return strings.ToLower(fmt.Sprintf("%s/%s:%s", addr, d.Project, d.Hash))
}
