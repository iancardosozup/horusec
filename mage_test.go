package horusec

import (
	"context"
	"errors"
	"fmt"
	"github.com/ZupIT/horusec/internal/utils/testutil"
	"github.com/google/go-github/v40/github"
	"github.com/magefile/mage/sh"
	"golang.org/x/mod/semver"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	Patch     = "patch"
	Minor     = "minor"
	Major     = "major"
	Alpha     = "alpha"
	Rc        = "rc"
	Beta      = "beta"
	Release   = "release"
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
	//		name:           "Should patch rc v0.0.1-rc.1 with success",
	//		version:        Patch,
	//		releaseType:    Rc,
	//		input:          V001RC1,
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
	version := Patch
	tag, err := getLatestReleaseTag("iancardosozup", "horusec", Beta)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tag)
	newTag, err := getNewReleaseTag(tag, version, Alpha)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(newTag)
}
func getLatestReleaseTag(org, project, releasetype string) (string, error) {
	if releasetype != Beta && releasetype != Alpha && releasetype != Rc && releasetype != Release {
		return "", fmt.Errorf("invalid release type, choose between %q %q %q %q", Alpha, Beta, Rc, Release)
	}
	ghClient := github.NewClient(nil)
	listOptions := &github.ListOptions{
		Page:    1,
		PerPage: 80,
	}
	tags, resp, err := ghClient.Repositories.ListTags(context.Background(), org, project, listOptions)
	if github.CheckResponse(resp.Response) != nil {
		return "", err
	}
	var latestTagName string
	for _, tag := range tags {
		if strings.Contains(tag.GetName(), releasetype) && releasetype != Release {
			latestTagName = tag.GetName()
			break
		}else if !strings.Contains(tag.GetName(), Beta) && !strings.Contains(tag.GetName(),Alpha) && !strings.Contains(tag.GetName(),Rc)  {
			latestTagName = tag.GetName()
			break
		}
	}
	return strings.ReplaceAll(latestTagName, `"`, ""), err
}
func getNewReleaseTag(currentTag, version, releaseType string) (string, error) {
	if !semver.IsValid(currentTag) {
		return "", errors.New("invalid current tag")
	}
	releaseType = strings.ToLower(releaseType)
	if releaseType != Beta && releaseType != Rc && releaseType != Release && releaseType != Alpha {
		return "", errors.New("invalid release type")
	}
	var releaseTag string
	currentTag = strings.ReplaceAll(currentTag, "v", "")
	versionSlice := strings.Split(currentTag, ".")

	nonOficialReleaseSlice, major, minor, patch, err := getSemverValues(currentTag, versionSlice)
	if err != nil {
		return "", err
	}
	releaseTag = fmt.Sprintf("%s%d.%d.%d", "v", major, minor, patch)
	switch version {
	case Patch:
		patch = patch + 1
		releaseTag = fmt.Sprintf("%s%d.%d.%d", "v", major, minor, patch)
		if releaseType == Rc {
			releaseTag, err = getNonOficialReleaseTagByCurrentTag(currentTag, Rc, nonOficialReleaseSlice, releaseTag, major, minor, patch)
		} else if releaseType == Beta {
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
		return "", fmt.Errorf("invalid version type, chose one: %q %q %q", Major, Minor, Patch)
	}
	err = os.Setenv("CLI_VERSION", releaseTag)
	if err != nil {
		return "", err
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

func TestB(t *testing.T) {
	newTag := Alpha

	if err := sh.RunV("git", "tag", "-f", newTag, "-m", "release "+newTag); err != nil {
		t.Error(err)
	}
	if err := sh.RunV("git", "push", "origin", "-f", newTag); err != nil {
		t.Error(err)
	}
	os.Setenv("GORELEASER_PREVIOUS_TAG", newTag)
	os.Setenv("CLI_VERSION", newTag)
	os.Setenv("CURRENT_DATE", time.Now().String())
	os.Setenv("COSIGN_KEY_LOCATION", filepath.Join(testutil.RootPath, "cosign.key"))
	os.Setenv("COSIGN_PWD", "123")
	sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")

	envs := map[string]string{
		"GORELEASER_PREVIOUS_TAG": os.Getenv("GORELEASER_PREVIOUS_TAG"),
		"CURRENT_DATE":            os.Getenv("CURRENT_DATE"),
		"COSIGN_KEY_LOCATION":     os.Getenv("COSIGN_KEY_LOCATION"),
		"COSIGN_PWD":              os.Getenv("COSIGN_PWD"),
	}
	var result []string
	for k, v := range envs {
		if v == "" {
			result = append(result, k)
		}
	}
}
