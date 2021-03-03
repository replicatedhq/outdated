package outdated

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/genuinetools/reg/registry"
	semver "github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

type VersionTag struct {
	Sort    int    `json:"sort"`
	Version string `json:"version"`
	Date    string `json:"date"`
}

type V1History struct {
	Created string `json:"created,omitempty"`
}

type SemverTagCollection []*semver.Version

func (c SemverTagCollection) Len() int {
	return len(c)
}

func (c SemverTagCollection) Less(i, j int) bool {
	return compareVersions(c[i], c[j]) < 0
}

func compareVersions(verI *semver.Version, verJ *semver.Version) int {
	if verI.LessThan(verJ) {
		return -1
	} else if verI.GreaterThan(verJ) {
		return 1
	}

	return 0
}

func (c SemverTagCollection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c SemverTagCollection) VersionsBehind(currentVersion *semver.Version) ([]*semver.Version, error) {
	cleaned, err := c.Unique()
	if err != nil {
		return []*semver.Version{}, errors.Wrap(err, "failed to deduplicate versions")
	}

	sortable := SemverTagCollection(cleaned)
	sort.Sort(sortable)

	for idx := range sortable {
		if sortable[idx].Original() == currentVersion.Original() {
			return sortable[idx:], nil
		}
	}

	return []*semver.Version{
		currentVersion,
	}, nil // /shrug
}

// Unique will create a new sorted slice with the same versions that have different tags removed.
// While this is valid in semver, it's used in docker images differently
// For example: redis:4-alpine and redis:4-debian are the same version
func (c SemverTagCollection) Unique() ([]*semver.Version, error) {
	unique := make(map[string]*semver.Version)

	for _, v := range c {
		var ver string
		var validSegments []int
		splitTag := strings.Split(v.Original(), ".")
		segments := v.Segments()

		if len(splitTag) == 1 {
			validSegments = []int{segments[0]}
		} else if len(splitTag) == 2 {
			validSegments = segments[0:2]
		} else if len(splitTag) == 3 {
			validSegments = segments
		}

		strSegments := []string{}
		for _, segment := range validSegments {
			strSegments = append(strSegments, strconv.Itoa(segment))
		}
		ver = strings.Join(strSegments, ".")

		if _, exists := unique[ver]; !exists {
			unique[ver] = v
		} else {
			// we want the shortest tag -
			// e.g. between redis:4-alpine and redis:4, we want redis:4
			if len(v.Original()) < len(unique[ver].Original()) {
				unique[ver] = v
			}
		}
	}

	result := make([]*semver.Version, 0, 0)
	for _, u := range unique {
		result = append(result, u)
	}

	sort.Sort(SemverTagCollection(result))

	return result, nil
}

// RemoveLeastSpecific given a sorted collection will remove the least specific version
func (c SemverTagCollection) RemoveLeastSpecific() []*semver.Version {
	if c.Len() == 0 {
		return []*semver.Version{}
	}

	cleanedVersions := []*semver.Version{c[0]}
	for i := 0; i < len(c)-1; i++ {
		j := i + 1
		iSegments := c[i].Segments()
		jSegments := c[j].Segments()

		isLessSpecific := true
		for idx, iSegment := range iSegments {
			if len(jSegments) < idx+1 {
				break
			}
			if iSegment > 0 && jSegments[idx] == 0 {
				break
			}
			if iSegment != jSegments[idx] {
				isLessSpecific = false
				break
			}
		}

		if !isLessSpecific {
			cleanedVersions = append(cleanedVersions, c[j])
		}
	}

	return cleanedVersions
}

func resolveTagDates(reg *registry.Registry, imageName string, sortedVersions []*semver.Version) ([]*VersionTag, error) {
	var wg sync.WaitGroup
	var mux sync.Mutex
	versionTags := make([]*VersionTag, 0)

	wg.Add(len(sortedVersions))
	for idx, version := range sortedVersions {
		versionFromTag := version.Original()
		versionTag := VersionTag{
			Sort:    idx,
			Version: versionFromTag,
		}

		go func(versionFromTag string) {
			date, err := getTagDate(reg, imageName, versionFromTag)
			if err == nil {
				versionTag.Date = date
			}

			mux.Lock()
			versionTags = append(versionTags, &versionTag)
			mux.Unlock()

			wg.Done()
		}(versionFromTag)

	}
	wg.Wait()

	return versionTags, nil
}

func getTagDate(reg *registry.Registry, imageName string, versionFromTag string) (string, error) {
	manifest, err := reg.ManifestV1(context.TODO(), imageName, versionFromTag)
	if err != nil {
		return "", errors.Wrap(err, "unable to get manifest from image")
	}
	for _, history := range manifest.History {
		v1History := V1History{}
		err := json.Unmarshal([]byte(history.V1Compatibility), &v1History)
		if err != nil {
			// if it doesn't fit...throw it away
			continue
		}
		if v1History.Created != "" {
			return v1History.Created, nil
		}
	}

	return "", errors.New("no dates found")
}
