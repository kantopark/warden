package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
	"github.com/pkg/errors"
)

// Tars a directory or file. If tarring is successful, returns the filepath
// of the tarred directory
func TarDir(source, target string) (string, error) {
	_, err := os.Stat(source)
	if err == os.ErrNotExist {
		return "", err
	}

	target = strings.TrimSpace(target)
	if StrIsEmptyOrWhitespace(target) {
		return "", errors.New("target name of tar file must not be empty")
	}
	// append .tar to target file name
	if filepath.Ext(target) != ".tar" {
		target = target + ".tar"
	}
	// replace invalid characters
	target = strings.Replace(target, ":", "-", -1)

	var files []string
	// walk through the source folder recursively and put in everything
	if err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		files = append(files, path)
		return err
	}); err != nil {
		return "", errors.Wrap(err, "error encountered when walking through source")
	}
	if err := archiver.Archive(files, target); err != nil {
		return "", errors.Wrap(err, "error encountered when tarring source")
	}

	return target, nil
}
