//go:build mage
// +build mage

// A comment on the package will be output when you list the targets of a
// magefile.
package main

import (
	"fmt"
	"github.com/ZupIT/horusec/internal/utils/testutil"
	"github.com/magefile/mage/sh"
	"os"
	"path/filepath"
	"time"
	// mage:import
	_ "github.com/ZupIT/horusec-devkit/pkg/utils/mageutils"
)

const (
	Patch       = "Patch"
	Minor       = "Minor"
	Major       = "Major"
	Alpha       = "alpha"
	Rc          = "rc"
	Beta        = "beta"
	TypeRelease = "release"
)

//GetCurrentDate execute "echo", `::set-output name=date::$(date "+%a %b %d %H:%M:%S %Y")`
func GetCurrentDate() error {
	if err := sh.RunV("echo", `::set-output name=date::$(date "+%a %b %d %H:%M:%S %Y")`); err != nil {
		return err
	}
	return nil
}

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

// tag returns the git tag for the current branch or "" if TypeRelease.
func tag() string {
	s, _ := sh.Output("git", "describe", "--tags")
	return s
}

// hash returns the git hash for the current repo or "" if TypeRelease.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

// ReleaseAlpha releases alpha version from main
func ReleaseAlpha() error {
	newTag := Alpha

	if err := sh.RunV("git", "tag", "-f", newTag, "-m", "release "+newTag); err != nil {
		return err
	}
	if err := sh.RunV("git", "push", "origin", "-f", newTag); err != nil {
		return err
	}

	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
}

//// ReleaseBeta releases beta version
//func ReleaseBeta(version string) (err error) {
//	mg.SerialDeps(verifyReleaseDeps)
//
//	newVersion, err := mage.NewVersion(version)
//	if err != nil {
//		return err
//	}
//	newTag := newVersion.NextBetaVersion
//	if err := sh.RunV("git", "tag", "-a", newTag, "-m", "release "+newTag); err != nil {
//		return err
//	}
//	if err := sh.RunV("git", "push", "origin", newTag); err != nil {
//		return err
//	}
//	defer func() {
//		if err != nil {
//			sh.RunV("git", "tag", "--delete", newTag)
//			sh.RunV("git", "push", "--delete", "origin", newTag)
//		}
//	}()
//	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
//}

//// ReleaseRc creates a release branch and a release candidate tag
//func ReleaseRc(version string) (err error) {
//	mg.Deps(verifyReleaseDeps)
//	var isReleaseBranchNew bool
//	tag, err := getLatestReleaseTag("iancardosozup", "horusec", Rc)
//	if err != nil {
//		return err
//	}
//	newTag, err := getNewReleaseTag(tag, version, Rc)
//	if err != nil {
//		return err
//	}
//	branchName := "rc/" + newTag[:4]
//	if err := sh.RunV("git", "checkout", branchName); err != nil {
//		isReleaseBranchNew = true
//		if err := sh.RunV("git", "checkout", "-b", branchName); err != nil {
//			return err
//		}
//	}
//
//	if err := sh.RunV("git", "tag", "-a", newTag, "-m", "release candidate"); err != nil {
//		return err
//	}
//	defer func() {
//		if err != nil {
//			sh.RunV("git", "checkout", "main")
//			sh.RunV("git", "tag", "--delete", newTag)
//			sh.RunV("git", "push", "--delete", "origin", newTag)
//			if isReleaseBranchNew {
//				sh.RunV("git", "branch", "-d", branchName)
//				sh.RunV("git", "push", "origin")
//			}
//		}
//	}()
//	if err := sh.RunV("git", "push", "origin", newTag); err != nil {
//		return err
//	}
//	if err := sh.RunV("git", "push", "origin", branchName); err != nil {
//		return err
//	}
//	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
//}

////Release creates a release by the release branch
//func Release(version string) (err error) {
//	mg.Deps(verifyReleaseDeps)
//	var isNewReleaseBranch bool
//	tag, err := getLatestReleaseTag("iancardosozup", "horusec", TypeRelease)
//	if err != nil {
//		return err
//	}
//	newTag, err := getNewReleaseTag(tag, version, TypeRelease)
//	if err != nil {
//		return err
//	}
//	rcBranchName := "rc/" + newTag[:4]
//	if err := sh.RunV("git", "checkout", rcBranchName); err != nil {
//		return err
//	}
//	releaseBranchName := "release/" + newTag[:4]
//	if err := sh.RunV("git", "checkout", releaseBranchName); err != nil {
//		isNewReleaseBranch = true
//		if err := sh.RunV("git", "checkout", "-b", releaseBranchName); err != nil {
//			return err
//		}
//	}
//	if err := sh.RunV("git", "tag", "-a", newTag, "-m", "release"); err != nil {
//		return err
//	}
//	if err := sh.RunV("git", "push", "origin", newTag); err != nil {
//		return err
//	}
//	defer func() {
//		if err != nil {
//			sh.RunV("git", "checkout", "main")
//			sh.RunV("git", "tag", "--delete", newTag)
//			sh.RunV("git", "push", "--delete", "origin", newTag)
//			if isNewReleaseBranch {
//				sh.RunV("git", "checkout", "main")
//				sh.RunV("git", "branch", "-d", releaseBranchName)
//				sh.RunV("git", "push", "origin")
//			}
//		}
//	}()
//	return sh.Run("goreleaser", "-f", filepath.Join(testutil.RootPath, "goreleaser.yml"), "--rm-dist")
//}

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
