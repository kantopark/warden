package deploy

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"

	"warden/store"
)

const (
	dockerPortMin = 40000
	dockerPortMax = 42673
)

var localIP string // A singleton string representing the IP address of the local machine.

func init() {
	if ip, err := getLocalIPAddress(); err != nil {
		log.Fatalln(err)
	} else {
		localIP = ip
	}
}

type dockerManager struct {
	routes *routeMap
	db     *store.Store
	ctx    context.Context
	cli    *client.Client
}

// Deploys an instance of the container on the Docker daemon
func (m *dockerManager) DeployInstance(d Deployment) error {
	port, err := findFreePort()
	if err != nil {
		return errors.Wrap(err, "could not find free port for deployment")
	}
	con, err := m.cli.ContainerCreate(
		m.ctx,
		&container.Config{
			Image:        d.ImageName(),
			ExposedPorts: nat.PortSet{nat.Port(port): struct{}{}},
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{nat.Port(port): {{HostIP: "localhost", HostPort: string(port)}}},
			AutoRemove:   true,
		},
		nil,
		d.ImageName())

	if err != nil {
		return errors.Wrap(err, "could not create instance")
	}
	if err := m.cli.ContainerStart(m.ctx, con.ID, types.ContainerStartOptions{}); err != nil {
		return errors.Wrap(err, "could not start instance")
	}
	m.routes.Set(d.Address(), fmt.Sprintf("%s:%d", localIP, port))

	return nil
}

// Stops the instance of the container on the Docker daemon. If deployment doesn't exist,
// nothing is done
func (m *dockerManager) StopInstance(d Deployment) error {
	ftr := filters.NewArgs()
	ftr.Add("name", d.ImageName())
	containers, _ := m.cli.ContainerList(
		m.ctx,
		types.ContainerListOptions{
			Filters: ftr,
			All:     true})

	for _, con := range containers {
		if err := m.cli.ContainerRemove(m.ctx, con.ID, types.ContainerRemoveOptions{
			Force: true,
		}); err != nil {
			return errors.Wrapf(err, "error removing container: %s", d.ImageName())
		}
	}
	m.routes.Delete(d.Address())
	return nil
}

func (m *dockerManager) Close() error {
	containers, err := m.cli.ContainerList(m.ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return errors.Wrap(err, "error getting all the containers for docker daemon shutdown")
	}
	// Remove all containers with public ports within range
	for _, con := range containers {
		for _, port := range con.Ports {
			if dockerPortMin <= port.PublicPort && port.PublicPort < dockerPortMax {
				if err := m.cli.ContainerRemove(m.ctx, con.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
					return errors.Wrapf(err, "error trying to remove container with image: %s ", con.ImageID)
				}
				break
			}
		}
	}

	// Close the docker client
	if err := m.cli.Close(); err != nil {
		return errors.Wrap(err, "error stopping docker deploy cli")
	}
	m.ctx.Done()
	return nil
}

func newDockerRunner() (*dockerManager, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "error creating new Docker Client")
	}

	db, err := store.NewStore()
	if err != nil {
		return nil, err
	}

	rm := newRouteMap()

	return &dockerManager{rm, db, ctx, cli}, nil
}

// Finds a free port between 40000 and 42367. This is usually only used when
// the runner is a base docker manager (local dev). A free port is found first
// before binding that port with docker run. Since this is meant for docker
// manager, the domain is assumed to always be localhost
func findFreePort() (int, error) {
	portChannel := make(chan int)
	for port := dockerPortMin; port < dockerPortMax; port++ {
		go func(testPort int) {
			conn, err := net.Dial("tcp", ":"+strconv.Itoa(testPort))
			if err != nil {
				portChannel <- testPort
				return
			}
			conn.Close()
		}(port)
	}

	select {
	case port := <-portChannel:
		return port, nil
	case <-time.After(20 * time.Second):
		return 0, errors.New("unable to find free ports between 40000 and 42367")
	}
}

// Gets the local IP address of the machine
func getLocalIPAddress() (string, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", errors.Wrap(err, "could not determine machine's IP address")
	}

	for _, a := range addresses {
		if n, ok := a.(*net.IPNet); ok && !n.IP.IsLoopback() && n.IP.To4() != nil && n.IP.IsGlobalUnicast() {
			return n.IP.String(), nil
		}
	}
	return "", errors.New("could not find global unicast address for machine. Use ifconfig (unix) or ipconfig (windows) to check if machine is connected to a network")
}
