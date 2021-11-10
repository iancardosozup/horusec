package horusec

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v40/github"
	"golang.org/x/mod/semver"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const (
	Patch     = "PATCH"
	Minor     = "MINOR"
	Major     = "MAJOR"
	Alpha     = "ALPHA"
	Rc        = "rc"
	Beta      = "beta"
	V001      = "v0.0.1"
	V001RC1   = "v0.0.1-rc.1"
	V001BETA1 = "v0.0.1-beta.1"
	V001RC2   = "v0.0.1-rc.2"
	V001BETA2 = "v0.0.1-beta.2"
	V002      = "v0.0.2"
	V002RC1   = "v0.0.2-rc.1"
	V002BETA1 = "v0.0.2-beta.1"
	V002RC2   = "v0.0.2-rc.2"
	V002BETA2 = "v0.0.2-beta.2"
	V003      = "v0.0.3"
	V003RC1   = "v0.0.3-rc.1"
	V003BETA1 = "v0.0.3-beta.1"
	V003RC2   = "v0.0.3-rc.2"
	V003BETA2 = "v0.0.3-beta.2"
	V010      = "v0.1.0"
	V010RC1   = "v0.1.0-rc.1"
	V010BETA1 = "v0.1.0-beta.1"
	V010RC2   = "v0.1.0-rc.2"
	V010BETA2 = "v0.1.0-beta.2"
	V020      = "v0.2.0"
	V020RC1   = "v0.2.0-rc.1"
	V020BETA1 = "v0.2.0-beta.1"
	V020RC2   = "v0.2.0-rc.2"
	V020BETA2 = "v0.2.0-beta.2"
	V030      = "v0.3.0"
	V030RC1   = "v0.3.0-rc.1"
	V030BETA1 = "v0.3.0-beta.1"
	V030RC2   = "v0.3.0-rc.2"
	V030BETA2 = "v0.3.0-beta.2"
	V100      = "v1.0.0"
	V100RC1   = "v1.0.0-rc.1"
	V100BETA1 = "v1.0.0-beta.1"
	V100RC2   = "v1.0.0-rc.2"
	V100BETA2 = "v1.0.0-beta.2"
	V200      = "v2.0.0"
	V200RC1   = "v2.0.0-rc.1"
	V200BETA1 = "v2.0.0-beta.1"
	V200RC2   = "v2.0.0-rc.2"
	V200BETA2 = "v2.0.0-beta.2"
)

func TestA(t *testing.T) {
	//testcases := []struct {
	//	name           string
	//	input          string
	//	releaseType    string
	//	version        string
	//	expectedOutput string
	//}{
	//	{
	//		name:           "Should patch none v0.0.1 with success",
	//		releaseType:    "",
	//		version:        Patch,
	//		input:          V001,
	//		expectedOutput: V002,
	//	},
	//	{
	//		name:           "Should patch rc v0.0.1 with success",
	//		version:        Patch,
	//		releaseType:    Rc,
	//		input:          V001,
	//		expectedOutput: V002RC1,
	//	},
	//	{
	//		name:        "Should patch rc v0.0.1-rc.1 with success",
	//		version:     Patch,
	//		releaseType: Rc,
	//		input:       V001RC1,
	//		//TODO: expectedOutput should be V002RC1 or V001RC2?
	//		expectedOutput: V001RC2,
	//	},
	//	{
	//		name:           "Should patch beta v0.0.1 with success",
	//		version:        Patch,
	//		releaseType:    Beta,
	//		input:          V001,
	//		expectedOutput: V002BETA1,
	//	}, {
	//		name:        "Should patch beta v0.0.1-beta.1 with success",
	//		version:     Patch,
	//		releaseType: Beta,
	//		input:       V001BETA1,
	//		//TODO: expectedOutput should be V002BETA1 or V001BETA2?
	//		expectedOutput: V001BETA2,
	//	}, {
	//		name:           "Should minor none v0.0.1 with success",
	//		version:        Minor,
	//		releaseType:    "",
	//		input:          V001,
	//		expectedOutput: V010,
	//	},
	//	{
	//		name:           "Should minor rc v0.0.1 with success",
	//		version:        Minor,
	//		releaseType:    Rc,
	//		input:          V001,
	//		expectedOutput: V010RC1,
	//	},
	//	{
	//		name:        "Should minor rc v0.0.1-rc.1 with success",
	//		version:     Minor,
	//		releaseType: Rc,
	//		input:       V010RC1,
	//		//TODO: expectedOutput should be V020RC1 or V010RC2?
	//		expectedOutput: V020RC2,
	//	},
	//	{
	//		name:           "Should minor beta v0.0.1 with success",
	//		version:        Minor,
	//		releaseType:    Beta,
	//		input:          V001,
	//		expectedOutput: V010BETA1,
	//	}, {
	//		name:           "Should Major none v0.0.1 with success",
	//		version:        Major,
	//		releaseType:    "",
	//		input:          V001,
	//		expectedOutput: V100,
	//	},
	//	{
	//		name:           "Should Major rc v0.0.1 with success",
	//		version:        Major,
	//		releaseType:    Rc,
	//		input:          V001,
	//		expectedOutput: V100RC1,
	//	},
	//	{
	//		name:           "Should Major beta v0.0.1 with success",
	//		version:        Major,
	//		releaseType:    Beta,
	//		input:          V001,
	//		expectedOutput: V100BETA1,
	//	},
	//	{
	//		name:        "Should Major rc v0.0.1-rc.1 with success",
	//		version:     Major,
	//		releaseType: Rc,
	//		input:       V001RC1,
	//		//TODO: expectedOutput should be V100RC1 or V100RC2?
	//		expectedOutput: V100RC2,
	//	},
	//	{
	//		name:           "Should Major v0.0.1-rc.1 with success",
	//		version:        Major,
	//		releaseType:    "",
	//		input:          V001RC1,
	//		expectedOutput: V100,
	//	},
	//	{
	//		name:           "Should minor v0.0.1-rc.1 with success",
	//		version:        Minor,
	//		releaseType:    "",
	//		input:          V001RC1,
	//		expectedOutput: V010,
	//	},
	//	{
	//		name:           "Should patch v0.0.1-rc.1 with success",
	//		version:        Patch,
	//		releaseType:    "",
	//		input:          V001RC1,
	//		expectedOutput: V002,
	//	},
	//}
	//for _, tt := range testcases {
	//	t.Run(tt.name, func(t *testing.T) {
	//
	//		resp, err := getNewReleaseTag(tt.input, tt.version, tt.releaseType)
	//		assert.NoError(t, err)
	//		assert.Equal(t, tt.expectedOutput, resp)
	//	})
	//
	//}
	version := Major
	tag, err := getLatestReleaseTag()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tag)
	newTag, err := getNewReleaseTag(tag, version, Beta)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(newTag)
}

func getLatestReleaseTag() (string, error) {
	ghClient := github.NewClient(nil)
	repo, resp, err := ghClient.Repositories.Get(context.Background(), "ZupIT", "horusec")
	if github.CheckResponse(resp.Response) != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodGet, strings.ReplaceAll(repo.GetReleasesURL(), `{/id}`, "/latest"), nil)
	var releases github.RepositoryRelease
	//date := time.Date(1999, time.May, 2, 1, 1, 1, 1, time.UTC)
	//latestRelease := github.RepositoryRelease{
	//	PublishedAt: &github.Timestamp{Time: date},
	//}

	resp, err = ghClient.Do(context.Background(), req, &releases)
	//if github.CheckResponse(resp.Response) != nil {
	//	return "", err
	//}
	//for _, release := range releases {
	//	if semver.IsValid(strings.ReplaceAll(github.Stringify(release.Name), `"`, "")) {
	//		latestRelease = release
	//		break
	//	}
	//}
	return strings.ReplaceAll(github.Stringify(releases.TagName), `"`, ""), err
}
func getNewReleaseTag(currentTag, version, releaseType string) (string, error) {
	if !semver.IsValid(currentTag) {
		return "", errors.New("invalid current tag")
	}
	var releaseTag string
	currentTag = strings.ReplaceAll(currentTag, "v", "")
	versionSlice := strings.Split(currentTag, ".")

	nonOficialReleaseSlice, major, minor, patch, err := getSemverValues(currentTag, versionSlice)
	if err != nil {
		return "", err
	}

	switch version {
	case Patch:
		patch = patch + 1
		releaseTag = fmt.Sprintf("%s%d.%d.%d", "v", major, minor, patch)
		if releaseType == Rc {
			patch = patch - 1
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Rc, nonOficialReleaseSlice, releaseTag, major, minor, patch)
		} else if releaseType == Beta {
			patch = patch - 1
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Beta, nonOficialReleaseSlice, releaseTag, major, minor, patch)
		}
	case Minor:
		minor = minor + 1
		releaseTag = fmt.Sprintf("%s%d.%d.%d", "v", major, minor, 0)
		if releaseType == Rc {
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Rc, nonOficialReleaseSlice, releaseTag, major, minor, 0)
		} else if releaseType == Beta {
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Beta, nonOficialReleaseSlice, releaseTag, major, minor, 0)
		}
	case Major:
		major = major + 1
		releaseTag = fmt.Sprintf("%s%d.%d.%d", "v", major, 0, 0)
		if releaseType == Rc {
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Rc, nonOficialReleaseSlice, releaseTag, major, 0, 0)
		} else if releaseType == Beta {
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Beta, nonOficialReleaseSlice, releaseTag, major, 0, 0)
		}
	case Alpha:
		releaseTag = "alpha"
	default:
		return "", errors.New("invalid release type")
	}

	return releaseTag, nil
}

func getSemverValues(currentTag string, versionSlice []string) (nonOficialReleaseSlice []string, major int, minor int, patch int, err error) {
	if strings.Contains(currentTag, "-") {
		nonOficialReleaseSlice = strings.Split(currentTag, "-")
		patchSlice := strings.Split(nonOficialReleaseSlice[0], ".")
		patch, err = strconv.Atoi(patchSlice[len(patchSlice)-1])
		if err != nil {
			return nil, 0, 0, 0, err
		}

	} else {
		patch, err = strconv.Atoi(versionSlice[2])
		if err != nil {
			return nil, 0, 0, 0, err
		}
	}
	major, err = strconv.Atoi(versionSlice[0])
	if err != nil {
		return nil, 0, 0, 0, err
	}
	minor, err = strconv.Atoi(versionSlice[1])
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return nonOficialReleaseSlice, major, minor, patch, nil
}

func getNonOficialReleaseTagByCurrentTag(currentTag, releaseType string, nonOficialReleaseSlice []string, releaseTag string, major int, minor int, patch int) (string, error) {
	if releaseType != Beta && releaseType != Rc {
		return "", fmt.Errorf("invalid release type: choose %q or %q", Beta, Rc)
	}
	if nonOficialReleaseSlice != nil {
		result := strings.Split(nonOficialReleaseSlice[1], ".")
		if result[len(result)-1] != "" {
			releaseTypeAttempt, err := strconv.Atoi(result[len(result)-1])
			if err != nil {
				return "", err
			}
			releaseTag = fmt.Sprintf("%s%d.%d.%d-%s.%d", "v", major, minor, patch, releaseType, releaseTypeAttempt+1)
		}
	} else {
		releaseTag = fmt.Sprintf("%s%s%s.%d", releaseTag, "-", releaseType, 1)
	}
	return releaseTag, nil
}
