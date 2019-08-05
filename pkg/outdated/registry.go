package outdated

import (
	"time"

	"github.com/docker/docker/api/types"
	"github.com/genuinetools/reg/registry"
	semver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

func initRegistryClient(hostname string) (*registry.Registry, error) {
	auth := types.AuthConfig{
		Username:      "",
		Password:      "",
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
		return nil, errors.Wrap(err, "list tags")
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
