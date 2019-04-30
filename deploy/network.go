package deploy

import (
	"net"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Finds a free port between 40000 and 42367. This is usually only used when
// the runner is a base docker manager (local dev). A free port is found first
// before binding that port with docker run. Since this is meant for docker
// manager, the domain is assumed to always be localhost
func findFreePort(min, max int) (int, error) {
	portChannel := make(chan int)
	for port := min; port < max; port++ {
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
