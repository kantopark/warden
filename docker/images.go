package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"warden/docker/templates"
	"warden/utils"
)

type ImageBuildOptions struct {
	Name     string
	GitURL   string
	Hash     string
	Username string
	Password string
	RunEnv   string
	Handler  string
	Alias    string
}

type templateDetails struct {
	Handler string
}

var box *templates.Box

func init() {
	_box, err := templates.NewBox()
	if err != nil {
		panic(err)
	}
	box = _box
}

// Builds the image. In nromal circumstances, you should run this as a go-routine
func (m *Manager) BuildImage(options ImageBuildOptions) error {
	// validations
	if utils.StrIsEmptyOrWhitespace(options.Handler) {
		return errors.New("handler must be specified")
	} else if utils.StrIsEmptyOrWhitespace(options.RunEnv) {
		return errors.New("RunEnv (runtime environment) must be specified")
	} else if utils.StrIsEmptyOrWhitespace(options.Username) {
		return errors.New("username must be specified")
	} else if utils.StrIsEmptyOrWhitespace(options.GitURL) {
		return errors.New("Repository (Git) url must be specified")
	} else if utils.StrIsEmptyOrWhitespace(options.Name) {
		return errors.New("project name must be specified")
	}

	// Cloning and checking out repository portion
	// Creating a temp folder to house the image build artifacts
	dir, err := ioutil.TempDir(os.TempDir(), options.Name+"-"+options.Hash)
	if err != nil {
		return errors.Wrap(err, "error creating temp dir for cloning when building image")
	}

	// Clone the repo into the temp folder
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: options.Username,
			Password: options.Password,
		},
		URL:      options.GitURL,
		Progress: os.Stdout,
	})

	if err != nil {
		return errors.Wrapf(err, "error cloning repo when building image. \n\tURL: %s. \n\tUsername: %s",
			options.GitURL, options.Username)
	}

	// Prepare the hash for the right checkout. If hash provided is an empty string or "latest",
	// Will checkout the latest commit
	options.Hash = strings.ToLower(strings.TrimSpace(options.Hash))
	if utils.StrIn(options.Hash, nil, "", "latest") {
		commits, err := repo.CommitObjects()
		if err != nil {
			return errors.Wrap(err, "error getting repo commits when building image")
		}
		commit, err := commits.Next()
		if err != nil {
			return errors.Wrap(err, "error getting latest commit when building image")
		}

		options.Hash = commit.Hash.String()
	}

	// Checkout the hash specified
	tree, err := repo.Worktree()
	if err != nil {
		return errors.Wrap(err, "error getting worktree when building image")
	}

	if err := tree.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(options.Hash),
	}); err != nil {
		return errors.Wrapf(err, "error checkout commit hash '%s' when building image", options.Hash)
	}

	// Create template Dockerfile in temp directory
	err = prepareDockerfileTemplate(dir, options.RunEnv, options.Handler)
	if err != nil {
		return errors.Wrap(err, "error when building image")
	}

	// Alias
	if utils.StrIsEmptyOrWhitespace(options.Alias) {
		options.Alias = "latest"
	} else {
		options.Alias = strings.ToLower(options.Alias)
	}

	tagName := fmt.Sprintf("%s:%s", options.Name, options.Alias)

	tarDir, err := utils.TarDir(dir, tagName, nil)
	if err != nil {
		return errors.Wrap(err, "error encountered when tarring payload for docker build context")
	}

	tarDir, _ = filepath.Abs(tarDir)
	tarfile, err := os.Open(tarDir)
	if err != nil {
		return errors.Wrap(err, "error encountered when reading tarfile")
	}
	defer tarfile.Close()

	// Preparing to build the image
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Build the image
	// TODO: add a repository
	resp, err := m.client.ImageBuild(ctx, tarfile, types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{tagName},
	})

	if err != nil {
		return errors.Wrap(err, "error encountered when building image")
	}
	defer resp.Body.Close()
	streamResponse(resp.Body)

	// remove temp directories and files
	os.RemoveAll(dir)
	os.Remove(tarDir)

	return nil
}

func (m *Manager) ListImages() ([]types.ImageSummary, error) {
	return m.client.ImageList(context.Background(), types.ImageListOptions{})
}

func prepareDockerfileTemplate(dir, env, handler string) error {
	file, err := os.Create(filepath.Join(dir, "Dockerfile"))
	if err != nil {
		return errors.Wrap(err, "error creating dockerfile template")
	}
	data := templateDetails{
		Handler: handler,
	}

	switch strings.ToLower(env) {
	case "python", "python3":
		tpl, err := box.GetTemplate("python")
		if err != nil {
			return err
		}

		if err := tpl.Execute(file, data); err != nil {
			return errors.Wrap(err, "error writing template dockerfile")
		}
	default:
		return errors.Errorf("Unknown runtime environment: %s", env)
	}

	return nil
}