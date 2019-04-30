package deploy

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"warden/utils"
)

// The payload information for the manager to determine where to
// send the function call to. It is sent from the client and redirected
// to the running instance with some modifications
type Payload struct {
	alias       string        // targeted project alias for request
	body        io.ReadCloser // payload from the user that will be sent to the instance
	headers     http.Header   // Request headers
	method      string        // request method
	project     string        // targeted project for request
	queryValues url.Values    // Query values
}

// Generates the address of the project given the project name and alias.
// Usually used as a key to inform the ingress controller of the endpoint
// for the instance that should serve the function call
func (p *Payload) Address() string {
	addr := p.project
	if p.alias != "" {
		addr += "/" + p.alias
	}

	return addr
}

// Executes the payload by running it through the Docker engine or Swarm/Kubernetes cluster
// The host domain of the cluster needs to be specified
func (p *Payload) Execute(host string) (*http.Response, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error running instance, failed creating cookie jar")
	}

	c := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           jar,
		Timeout:       5 * time.Minute,
	}

	req, err := http.NewRequest(p.method, p.getUrl(host), p.body)
	if err != nil {
		return nil, err
	}
	req.Header = p.headers
	return c.Do(req)
}

// Constructs the url to send the payload to. This url is the url to the
// instance running in the Docker engine or Swarm/Kubernetes instance.
func (p *Payload) getUrl(host string) string {
	var addr strings.Builder

	addr.WriteString(host)
	if len(p.queryValues) > 0 {
		addr.WriteString("?")
		for k, v := range p.queryValues {
			addr.WriteString(k + "=" + strings.Join(v, ","))
		}
	}
	return addr.String()
}

// Creates a new payload object from the client's request
func NewPayload(r *http.Request) (*Payload, error) {
	p := &Payload{
		project:     utils.StrLowerTrim(chi.URLParam(r, "project")),
		alias:       utils.StrLowerTrim(chi.URLParam(r, "alias")),
		headers:     r.Header,
		queryValues: r.URL.Query(),
	}

	if p.alias == "latest" {
		p.alias = ""
	}
	switch strings.ToUpper(r.Method) {
	case "GET":
		p.method = "GET"
		p.body = nil
	case "POST":
		p.method = "POST"
		p.body = r.Body
	default:
		return nil, errors.Errorf("Only GET and POST method are allowed")
	}

	return p, nil
}
