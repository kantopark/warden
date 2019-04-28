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
func findFreePort() (int, error) {
	portChannel := make(chan int)
	for port := 40000; port < 42500; port++ {
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
