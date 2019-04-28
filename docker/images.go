package docker

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"warden/docker/templates"
	"warden/utils"
)

// Standardized options for building function images. One should ensure that
// the username and password / access token provided has the authority to
// clone the repository
type ImageBuildOptions struct {
	Name     string // Project Name
	GitURL   string // URL of git repository
	Hash     string // Commit hash of git url to build
	Username string // Username used to clone repository.
	Password string // Either the password for the username account or preferably the access token
	RunEnv   string // Run time environment. i.e. Python
	Handler  string // Handler specifies the file and function that serves as the entrypoint. i.e. main.entry_func
	Alias    string // Alias for the function run
	buildId  string // Internal ID used to track whether image is getting built
}

// ImagePullOptions holds information to pull images.
type ImagePullOptions struct {
	types.ImagePullOptions
	UseDockerHub bool // If true, pulls image from docker hub, otherwise, pulls from private registry. Default false
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

// Builds the image specified in the ImageBuildOptions. In normal circumstances,
// you should run this as a go routine.
func (c *Client) BuildImage(options ImageBuildOptions) error {
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
	options.Name = strings.ToLower(options.Name)
	options.Hash = strings.ToLower(options.Hash)
	options.buildId = strings.ToLower(fmt.Sprintf("%s-%s", options.GitURL, options.Hash))

	if !utils.StrIsEmptyOrWhitespace(options.Hash) {
		if hasTag, err := c.hubHasImage(options.Name, options.Hash); err != nil {
			log.Println(err)
		} else if hasTag {
			log.Printf("image '%s:%s' already exists", options.Name, options.Hash)
			return nil // image exists in repository. Skip
		}
	}

	// Should probably set a lock here, but I don't foresee that kind of traffic
	res := c.redis.Get(options.buildId)
	if res.Err() != redis.Nil {
		// there is a similar image current building. Skip
		log.Printf("image '%s:%s' already building", options.Name, options.Hash)
		return nil
	}

	// default build time: 10 minutes
	c.redis.Set(options.buildId, fmt.Sprintf("Building image: %s", options.buildId), 10*time.Minute)

	// build image
	go func() {
		err := c.buildImage(options)
		if err != nil {
			c.redis.Set(
				options.buildId,
				fmt.Sprintf("Image build '%s' resulted in error. Check build again. If local build succeeded, it may mean that image build took more than 10 minutes (timeout error)", options.buildId),
				24*time.Hour)
			log.Println(err)
		}
	}()

	return nil
}

func (c *Client) buildImage(options ImageBuildOptions) error {
	// Cloning and checking out repository portion
	// Creating a temp folder to house the image build artifacts
	dir, err := ioutil.TempDir(os.TempDir(), options.Name+"-"+options.Hash)
	defer os.RemoveAll(dir)
	if err != nil {
		return errors.Wrap(err, "error creating temp dir for cloning when building image")
	}

	// Clone the repo into the temp folder
	user := options.Username
	pw := options.Password
	if utils.StrIsEmptyOrWhitespace(options.Password) {
		// if empty password, assume that no authorization needed to clone repo
		// If user is not empty but password is empty, will have error cloning anyway
		user, pw = "", ""
	}

	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth: &http.BasicAuth{
			Username: user,
			Password: pw,
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
	commits, err := repo.CommitObjects()
	if err != nil {
		return errors.Wrap(err, "error getting repo commits")
	}

	if utils.StrIn(options.Hash, nil, "", "latest") {
		commit, err := commits.Next()
		if err != nil {
			return errors.Wrap(err, "error getting latest commit when building image")
		}
		options.Hash = commit.Hash.String()
	} else {
		for {
			c, err := commits.Next()
			if err == io.EOF {
				return errors.Wrapf(err, "could not find commit prefixed with hash '%s'", options.Hash)
			} else if err != nil {
				return errors.Wrap(err, "error reading commit history")
			}
			if strings.HasPrefix(c.Hash.String(), options.Hash) {
				options.Hash = c.Hash.String()
				break
			}
		}
	}

	if hasTag, err := c.hubHasImage(options.Name, options.Hash); err != nil {
		log.Println(err)
	} else if hasTag {
		return nil // image exists, skip
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

	tagName := formRegistryTag(options.Name, options.Hash)

	tarDir, err := utils.TarDir(dir, tagName, &utils.TarDirOption{RemoveIfExist: true})
	if err != nil {
		return errors.Wrap(err, "error encountered when tarring payload for docker build context")
	}

	tarDir, _ = filepath.Abs(tarDir)
	defer os.Remove(tarDir)
	tarfile, err := os.Open(tarDir)
	if err != nil {
		return errors.Wrap(err, "error encountered when reading tarfile")
	}
	defer tarfile.Close()

	// Preparing to build the image
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Build the image
	if resp, err := c.cli.ImageBuild(ctx, tarfile, types.ImageBuildOptions{
		SuppressOutput: false,
		Remove:         true,
		ForceRemove:    true,
		PullParent:     true,
		Tags:           []string{tagName},
	}); err != nil {
		return errors.Wrap(err, "error encountered when building image")
	} else {
		defer resp.Body.Close()
		streamResponse(resp.Body)
	}

	// Push image to local registry.
	if resp, err := c.cli.ImagePush(
		c.ctx,
		tagName,
		types.ImagePushOptions{
			RegistryAuth: `Base64Encode{"username":username,"password":password}`,
		},
	); err != nil {
		return errors.Wrap(err, "error encountered when pushing image to (private) registry")
	} else {
		defer resp.Close()
		streamResponse(resp)
	}

	c.redis.Del(options.buildId)
	return nil
}

func (c *Client) ListImages() ([]types.ImageSummary, error) {
	return c.cli.ImageList(context.Background(), types.ImageListOptions{})
}

// Returns the first image specified by the regex name string. If multiple images
// are matched by the name, returns the first image
func (c *Client) FindImageByName(regexName string) (*types.ImageSummary, error) {
	ftr := filters.NewArgs()
	ftr.Add("reference", regexName)

	images, err := c.cli.ImageList(c.ctx, types.ImageListOptions{
		All:     false,
		Filters: ftr,
	})

	if err != nil {
		return nil, errors.Wrap(err, "error looking for image")
	}
	if len(images) == 0 {
		return nil, errors.New("could not find image specified")
	}

	return &images[0], nil
}

func (c *Client) PullImage(name string, options *ImagePullOptions) error {
	if options == nil {
		options = &ImagePullOptions{
			types.ImagePullOptions{},
			false,
		}
	}

	var repo string
	if options.UseDockerHub {
		repo = "docker.io/library"
	} else {
		repo := fmt.Sprintf("%s://%s", viper.GetString("registry.protocol"), viper.GetString("registry.domain"))
		port := viper.GetInt("registry.port")
		if port != 0 && port != 80 && port != 443 {
			repo = fmt.Sprintf("%s:%d", repo, port)
		}
	}

	resp, err := c.cli.ImagePull(c.ctx, fmt.Sprintf("%s/%s", repo, name), options.ImagePullOptions)
	if err != nil {
		return err
	}
	defer resp.Close()
	streamResponse(resp)

	return nil
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
