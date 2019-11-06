package outdated

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/genuinetools/reg/registry"
	semver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

type DockerConfig struct {
	Auths map[string]DockerAuth `json:"auths"`
}

type DockerAuth struct {
	Auth string `json:"auth"`
}

func initRegistryClient(hostname string) (*registry.Registry, error) {
	if hostname == "docker.io" {
		hostname = "index.docker.io"
	}

	username := ""
	password := ""

	if _, err := os.Stat(path.Join(homeDir(), ".docker", "config.json")); err == nil {
		data, err := ioutil.ReadFile(path.Join(homeDir(), ".docker", "config.json"))
		if err != nil {
			return nil, errors.Wrap(err, "failed to read docker config")
		}

		dockerConfig := DockerConfig{}
		if err := json.Unmarshal(data, &dockerConfig); err != nil {
			return nil, errors.Wrap(err, "failed to parse docker config")
		}

		for host, auth := range dockerConfig.Auths {
			useAuth := false
			if strings.Contains(host, "index.docker.io") && hostname == "index.docker.io" {
				useAuth = true
			}
			if host == hostname {
				useAuth = true
			}

			if useAuth {
				decoded, err := base64.StdEncoding.DecodeString(auth.Auth)
				if err != nil {
					return nil, errors.Wrap(err, "failed to decode docker auth")
				}

				usernameAndPassword := strings.Split(string(decoded), ":")
				if len(usernameAndPassword) == 2 {
					username = usernameAndPassword[0]
					password = usernameAndPassword[1]
				}
			}
		}
	}

	auth := types.AuthConfig{
		Username:      username,
		Password:      password,
		ServerAddress: hostname,
	}

	reg, err := registry.New(auth, registry.Opt{
		SkipPing: true,
		Timeout:  time.Duration(time.Second * 5),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create registry client")
	}

	return reg, nil
}

func fetchTags(reg *registry.Registry, imageName string) ([]string, error) {
	tags, err := reg.Tags(imageName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list tags")
	}

	return tags, nil
}

func parseTags(tags []string) ([]*semver.Version, []string) {
	semverTags := make([]*semver.Version, 0, 0)
	nonSemverTags := make([]string, 0, 0)

	for _, tag := range tags {
		v, err := semver.NewVersion(tag)
		if err != nil {
			nonSemverTags = append(nonSemverTags, tag)
		} else {
			semverTags = append(semverTags, v)
		}
	}

	return semverTags, nonSemverTags
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
