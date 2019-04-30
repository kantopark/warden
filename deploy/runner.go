package deploy

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"warden/utils"
)

var once sync.Once
var manager Manager
var managerError error

type Manager interface {
	Close() error
	DeployInstance(d Deployment) error
	StopInstance(d Deployment) error
}

// Returns a deployment manager. Manager is a singleton object. Manager is
// used as a control layer to handle deployment of functions (instances)
// for the different deployment runtime (Docker, Swarm or Kubernetes). Note
// that Docker should only be used for local testing while the choice of
// Swarm or Kubernetes depends on your preference
func NewManager() (Manager, error) {
	once.Do(func() {
		r := utils.StrLowerTrim(viper.GetString("deploy.type"))
		switch r {
		case "docker":
			manager, managerError = newDockerRunner()
		case "swarm":
			managerError = errors.New("Swarm manager is not yet implemented")
		case "kubernetes", "k8s":
			managerError = errors.New("Kubernetes manager is not yet implemented")
		default:
			managerError = errors.Errorf("Unknown deploy type: %s", r)
		}
	})
	return manager, managerError
}

type Deployment struct {
	Alias      string
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

	d.Alias = utils.StrLowerTrim(d.Alias)
	// Latest alias also becomes default route
	if d.Alias == "latest" {
		d.Alias = ""
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

// Gets the Image name for the deployment
func (d *Deployment) ImageName() string {
	addr := viper.GetString("registry.domain")
	port := viper.GetInt("registry.port")
	if port != 0 && port != 80 && port != 443 {
		addr = fmt.Sprintf("%s:%d", addr, port)
	}

	return utils.StrLowerTrim(fmt.Sprintf("%s/%s:%s", addr, d.Project, d.Hash))
}

// Gets the tail address (without the domain) for the deployment.
func (d *Deployment) Address() string {
	route := d.Project
	if d.Alias != "" {
		route += "/" + d.Alias
	}
	return utils.StrLowerTrim(route)
}
