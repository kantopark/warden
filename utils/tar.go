package utils

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Tars a directory or file. If tarring is successful, returns the filepath
// of the tarred directory
func TarDirectory(source, target string) (string, error) {
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

	// create file
	writer, err := os.Create(target)
	if err != nil {
		return "", errors.Wrap(err, "error creating tarfile")
	}
	defer writer.Close()

	// set file as a tar writer
	tarfile := tar.NewWriter(writer)
	defer tarfile.Close()

	// walk through the source folder recursively and put in everything
	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		// Skip the git folder
		if info.Name() == ".git" && info.IsDir() {
			return nil
		}

		// return on any error
		if err != nil {
			return err
		}

		// return on non-regular files
		if !info.Mode().IsRegular() {
			return nil
		}

		// set file headers for each file in the tar folder
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return errors.Wrap(err, "error setting tar header")
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(path, source, "", -1), string(filepath.Separator))
		if err := tarfile.WriteHeader(header); err != nil {
			return err
		}

		// copy file into tar folder
		file, err := os.Open(path)
		if err != nil {
			return errors.Wrapf(err, "error open file at: %s", path)
		}
		defer file.Close()

		_, err = io.Copy(tarfile, file)
		return err
	})
	if err != nil {
		return "", errors.Wrap(err, "error tarring source")
	}

	return target, nil
}
