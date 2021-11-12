//go:build mage
// +build mage

// A comment on the package will be output when you list the targets of a
// magefile.
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ZupIT/horusec/internal/utils/testutil"
	"github.com/google/go-github/v40/github"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"golang.org/x/mod/semver"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	Patch = "patch"
	Minor = "minor"
	Major = "major"
	Alpha = "alpha"
	Rc    = "rc"
	Beta  = "beta"
	None  = ""
)

// Runs go mod download and then installs the binary.
func Build() error {
	sh.RunV("go", "mod", "download")
	return sh.RunV("go", "build", "-o", "project", "-ldflags="+ldflags(), "github.com/ZupIT/horusec")
}

func ldflags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := hash()
	tag := tag()
	if tag == "" {
		tag = "dev"
	}
	return fmt.Sprintf(`-X "github.com/Mattel/project/proj.timestamp=%s" `+
		`-X "github.com/Mattel/project/proj.commitHash=%s" `+
		`-X "github.com/Mattel/project/proj.gitTag=%s"`, timestamp, hash, tag)
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// ReleaseAlpha releases alpha version from main
func ReleaseAlpha() error {
	mg.Deps(verifyReleaseDeps)
	newTag := Alpha

	if err := sh.RunV("git", "tag", "-f", newTag, "-m", "release "+newTag); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", "-f", newTag); err != nil {
		return err
	}
	os.Setenv("GORELEASER_PREVIOUS_TAG", newTag)
	os.Setenv("CLI_VERSION", newTag)
	os.Setenv("CURRENT_DATE", time.Now().String())
	os.Setenv("COSIGN_KEY_LOCATION", filepath.Join(testutil.RootPath, "cosign.key"))
	os.Setenv("COSIGN_PWD", "123")
	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
}

// ReleaseBeta releases beta version
func ReleaseBeta(version string) (err error) {
	mg.Deps(verifyReleaseDeps)
	tag, err := getLatestReleaseTag()

	if err != nil {
		return err
	}
	newTag, err := getNewReleaseTag(tag, version, Beta)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sh.RunV("git", "tag", "--delete", newTag)
			sh.RunV("git", "push", "--delete", "origin", newTag)
		}
	}()
	if err := sh.RunV("git", "tag", "-a", newTag, "-m", "release "+newTag); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", newTag); err != nil {
		return err
	}
	os.Setenv("GORELEASER_PREVIOUS_TAG", newTag)
	os.Setenv("CLI_VERSION", newTag)
	os.Setenv("CURRENT_DATE", time.Now().String())
	os.Setenv("COSIGN_KEY_LOCATION", filepath.Join(testutil.RootPath, "cosign.key"))
	os.Setenv("COSIGN_PWD", "123")
	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
}

// ReleaseRc creates a release branch and a release candidate tag
func ReleaseRc(version string) (err error) {
	mg.Deps(verifyReleaseDeps)
	tag, err := getLatestReleaseTag()
	if err != nil {
		return err
	}
	newTag, err := getNewReleaseTag(tag, version, Rc)
	if err != nil {
		return err
	}
	versionSlice := strings.Split(newTag, "-")[0]

	branchName := "release/" + versionSlice[:len(versionSlice)-2]
	if err := sh.RunV("git", "checkout", "-b", branchName); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sh.RunV("git", "checkout", "main")
			sh.RunV("git", "branch", "-d", branchName)
			sh.RunV("git", "tag", "--delete", newTag)
			sh.RunV("git", "push", "--delete", "origin", newTag)
		}
	}()

	if err := sh.RunV("git", "tag", "-a", newTag, "-m", "release candidate"); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", newTag); err != nil {
		return err
	}
	os.Setenv("GORELEASER_PREVIOUS_TAG", newTag)
	os.Setenv("CLI_VERSION", newTag)
	os.Setenv("GITHUB_TOKEN", os.Getenv("HORUSEC_PUSH_TOKEN"))
	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
}

func verifyReleaseDeps() error {
	err := hasAllNecessaryEnvs()
	if err != nil {
		return err
	}
	return hasGoreleaser()
}
func hasGoreleaser() error {
	return sh.Run("goreleaser", "-h")
}
func hasAllNecessaryEnvs() error {
	envs := map[string]string{
		"HORUSEC_PUSH_TOKEN":  os.Getenv("HORUSEC_PUSH_TOKEN"),
		"CURRENT_DATE":        os.Getenv("CURRENT_DATE"),
		"COSIGN_KEY_LOCATION": os.Getenv("COSIGN_KEY_LOCATION"),
		"COSIGN_PWD":          os.Getenv("COSIGN_PWD"),
	}
	var result []string
	for k, v := range envs {
		if v == "" {
			result = append(result, k)
		}
	}
	if len(result) != 0 {
		return fmt.Errorf("Missing some env var: %v", result)
	}
	return nil
}
func getLatestReleaseTag() (string, error) {
	ghClient := github.NewClient(nil)
	repo, resp, err := ghClient.Repositories.Get(context.Background(), "iancardosozup", "horusec")
	if github.CheckResponse(resp.Response) != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodGet, strings.ReplaceAll(repo.GetReleasesURL(), `{/id}`, "/latest"), nil)
	var release github.RepositoryRelease
	resp, err = ghClient.Do(context.Background(), req, &release)
	if github.CheckResponse(resp.Response) != nil {
		return "", err
	}
	return strings.ReplaceAll(github.Stringify(release.TagName), `"`, ""), err
}
func getNewReleaseTag(currentTag, version, releaseType string) (string, error) {
	if !semver.IsValid(currentTag) {
		return "", errors.New("invalid current tag")
	}
	releaseType = strings.ToLower(releaseType)
	if releaseType != Beta && releaseType != Rc && releaseType != None && releaseType != Alpha {
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
