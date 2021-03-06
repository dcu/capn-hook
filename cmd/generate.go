// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"github.com/dcu/capn-hook/core"
	"github.com/mattn/go-zglob"
	"github.com/spf13/cobra"
)

// generatorCmd represents the generator command
var generatorCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates a default config for your project.",
	Long:    `Automatically detects if your project is golang or ruby and generates a default manifest for it.`,
	Aliases: []string{"g"},
	Run: func(cmd *cobra.Command, args []string) {
		var manifest *core.Manifest
		if countFiles("**/*.go") > 0 {
			manifest = core.DefaultGolangManifest()
		} else if countFiles("**/*.rb") > 0 {
			manifest = core.DefaultRubyManifest()
		} else if files := findFiles("build.gradle"); len(files) > 0 { // gradle project
			manifest = core.DefaultAndroidManifest()
		} else {
			manifest = core.DefaultManifest()
		}

		fmt.Printf("Writing default config file to: %s\n", core.DefaultManifestFileName)
		manifest.WriteFile(core.DefaultManifestFileName)
	},
}

func findFiles(pattern string) []string {
	files, err := zglob.Glob(pattern)
	if err != nil {
		return []string{}
	}

	return files
}

func countFiles(pattern string) int {
	return len(findFiles(pattern))
}

func init() {
	RootCmd.AddCommand(generatorCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generatorCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generatorCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
