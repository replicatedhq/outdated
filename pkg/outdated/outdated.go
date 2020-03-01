package outdated

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/genuinetools/reg/registry"
	semver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

type Outdated struct {
}

type CheckResult struct {
	IsAccessible   bool
	LatestVersion  string
	VersionsBehind int64
	CheckError     string
	Path           string
}

func (o Outdated) ParseImage(image string, pullableImage string) (*CheckResult, error) {
	hostname, imageName, tag, err := ParseImageName(image)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse image name")
	}

	reg, err := initRegistryClient(hostname)
	if err != nil {
		return &CheckResult{
			IsAccessible:   false,
			LatestVersion:  "",
			VersionsBehind: -1,
			CheckError:     fmt.Sprintf("Cannot access registry: %s", err.Error()),
		}, nil
	}

	tags, err := fetchTags(reg, imageName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch image tags")
	}

	semverTags, nonSemverTags, err := parseTags(tags)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse tags")
	}

	detectedSemver, err := semver.NewVersion(tag)
	if err != nil {
		return o.parseNonSemverImage(reg, imageName, tag, nonSemverTags)
	}

	// From here on, we can assume that we are on a semver tag
	semverTags = append(semverTags, detectedSemver)
	collection := SemverTagCollection(semverTags)

	versionsBehind, err := collection.VersionsBehind(detectedSemver)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate versions behind")
	}
	trueVersionsBehind := SemverTagCollection(versionsBehind).RemoveLeastSpecific()

	behind := len(trueVersionsBehind) - 1

	checkResult := CheckResult{
		IsAccessible: true,
	}
	checkResult.VersionsBehind = int64(behind)

	versionPaths, err := resolveTagDates(reg, imageName, trueVersionsBehind)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve tag dates")
	}
	path, err := json.Marshal(versionPaths)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal version path")
	}
	checkResult.Path = string(path)

	checkResult.LatestVersion = trueVersionsBehind[len(trueVersionsBehind)-1].String()

	return &checkResult, nil
}

func (o *Outdated) parseNonSemverImage(reg *registry.Registry, imageName string, tag string, nonSemverTags []string) (*CheckResult, error) {
	laterDates := []string{}
	tagDate, err := getTagDate(reg, imageName, tag)
	if err != nil {
		return &CheckResult{
			IsAccessible:   true,
			LatestVersion:  tag,
			VersionsBehind: -1,
			CheckError:     "Unable to determine date from current tag",
		}, nil
	}
	myDate, err := time.Parse(time.RFC3339Nano, tagDate)
	if err != nil {
		return nil, err
	}

	for _, nonSemverTag := range nonSemverTags {
		otherDate, err := getTagDate(reg, imageName, nonSemverTag)
		if err != nil {
			continue
		}

		o, err := time.Parse(time.RFC3339Nano, otherDate)
		if err != nil {
			continue
		}
		if o.After(myDate) {
			laterDates = append(laterDates, otherDate)
		}
	}

	behind := int64(len(laterDates))
	return &CheckResult{
		IsAccessible:   true,
		LatestVersion:  tag,
		VersionsBehind: behind,
		CheckError:     "",
	}, nil
}
