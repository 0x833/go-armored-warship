// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// Version set by the build
var (
	Branch string
	SHA    string
	SemVer string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display the current version of the game",
	Long:  `display the current version of the game`,
	Run: func(cmd *cobra.Command, args []string) {
		versionInfo()
	},
}

func versionInfo() {
	repeatCount := 40
	fmt.Println(strings.Repeat("=", repeatCount))
	fmt.Printf("SemVer: %s\nBranch: %s\nSHA: %s\n", SemVer, Branch, SHA)
	fmt.Println(strings.Repeat("=", repeatCount))
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
