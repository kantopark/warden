package deploy

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"

	"warden/utils"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

// The payload information for the manager to determine where to
// send the function call to. It is sent from the client and redirected
// to the running instance with some modifications
type Payload struct {
	Alias       string        // targeted project alias for request
	Body        io.ReadCloser // payload from the user that will be sent to the instance
	Method      Method        // request method
	Project     string        // targeted project for request
	Headers     http.Header   // Request headers
	queryValues url.Values    // Query values
}

// Generates the address of the project given the project name and alias.
// Usually used as a key to inform the ingress controller of the endpoint
// for the instance that should serve the function call
func (p *Payload) Address() string {
	var addr strings.Builder

	addr.WriteString(p.Project)

	if p.Alias != "" {
		addr.WriteString("/" + p.Alias)
	}

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
		Project:     utils.StrLowerTrim(chi.URLParam(r, "project")),
		Alias:       utils.StrLowerTrim(chi.URLParam(r, "alias")),
		Body:        r.Body,
		Headers:     r.Header,
		queryValues: r.URL.Query(),
	}

	if p.Alias == "latest" {
		p.Alias = ""
	}
	switch strings.ToUpper(r.Method) {
	case "GET":
		p.Method = GET
	case "POST":
		p.Method = POST
	default:
		return nil, errors.Errorf("Only GET and POST method are allowed")
	}

	return p, nil
}
