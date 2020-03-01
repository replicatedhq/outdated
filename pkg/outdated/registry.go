package outdated

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/genuinetools/reg/registry"
	semver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

const (
	// SemverOutlierMajorVersionThreshold defines the number of major versions that must be skipped before
	// the next version is considered an outlier
	// setting this to 2 allows only 1 major version to be skipped
	SemverOutlierMajorVersionThreshold = 2
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

func parseTags(tags []string) ([]*semver.Version, []string, error) {
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

	// some semver tags might be outliers and should be treated as non-semver tags
	// For more info, see https://github.com/replicatedhq/outdated/issues/19
	outlierSemver, remainingSemver, err := splitOutlierSemvers(semverTags)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to split outliers")
	}

	for _, outlier := range outlierSemver {
		nonSemverTags = append(nonSemverTags, outlier.String())
	}

	return remainingSemver, nonSemverTags, nil
}

func splitOutlierSemvers(allSemverTags []*semver.Version) ([]*semver.Version, []*semver.Version, error) {
	if len(allSemverTags) == 0 {
		return []*semver.Version{}, []*semver.Version{}, nil
	}

	sortable := SemverTagCollection(allSemverTags)
	sort.Sort(sortable)

	outliers := []*semver.Version{}
	remaining := []*semver.Version{}

	lastVersion := allSemverTags[0]
	isInOutlier := false
	for _, v := range allSemverTags {
		if v.Segments()[0]-lastVersion.Segments()[0] > SemverOutlierMajorVersionThreshold {
			isInOutlier = true
		}

		if isInOutlier {
			outliers = append(outliers, v)
		} else {
			remaining = append(remaining, v)
		}

		lastVersion = v
	}

	return outliers, remaining, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}
