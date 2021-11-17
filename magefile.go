//go:build mage
// +build mage

// A comment on the package will be output when you list the targets of a
// magefile.
package main

import (
	"github.com/magefile/mage/sh"
	// mage:import
	_ "github.com/ZupIT/horusec-devkit/pkg/utils/mageutils"
)

// GetCurrentDate execute "echo", `::set-output name=date::$(date "+%a %b %d %H:%M:%S %Y")`
func GetCurrentDate() error {
	if err := sh.RunV("echo", `::set-output name=date::$(date "+%a %b %d %H:%M:%S %Y")`); err != nil {
		return err
	}
	return nil
}
