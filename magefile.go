// Copyright 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build mage

// A comment on the package will be output when you list the targets of a
// magefile.
package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	// mage:import
	_ "github.com/ZupIT/horusec-devkit/pkg/utils/mageutils"
)

// GetCurrentDate execute "echo", `::set-output name=date::$(date "+%a %b %d %H:%M:%S %Y")`
func GetCurrentDate() error {
	return sh.RunV("echo", `::set-output name=date::$(date "+%a %b %d %H:%M:%S %Y")`)
}
func SetPotatoVariable() error {
	//stdOut := os.Stdout
	//_, err := stdOut.WriteString("echo potato=fried_potato >> $GITHUB_ENV")
	//if err != nil {
	//	return err
	//}
	//exec.Command("echo potato=fried_potato >> $GITHUB_ENV")
	println("echo potato=fried_potato >> $GITHUB_ENV ")
	fmt.Println("echo potato=fried_potato >> $GITHUB_ENV ")
	return sh.Run("echo", "potato=fried_potato", ">>", "$GITHUB_ENV")
}
