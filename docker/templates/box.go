package templates

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type Box struct {
	templates map[string]*template.Template
}

func NewBox() (*Box, error) {
	_, dir, _, _ := runtime.Caller(0)
	dir = filepath.Dir(dir)

	_map := make(map[string]*template.Template)

	// get all files in directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "error encountered when creating template box")
	}

	for _, file := range files {
		name := file.Name()
		if file.IsDir() || filepath.Ext(name) != ".dtpl" {
			continue // skip if not template file
		}

		fp := filepath.Join(dir, name)
		// save file as template object
		tpl, err := template.ParseFiles(fp)
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing template file: %s", name)
		}
		name = strings.TrimSuffix(name, filepath.Ext(name))
		_map[name] = tpl
	}

	return &Box{
		templates: _map,
	}, nil
}

func (b *Box) GetTemplate(name string) (*template.Template, error) {
	tpl, ok := b.templates[name]
	if !ok {
		return nil, errors.Errorf("No template with key: %s", name)
	}
	return tpl, nil
}
