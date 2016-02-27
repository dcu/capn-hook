// Copyright © 2016 David Cuadrado
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
	"io/ioutil"
	"path/filepath"

	"code.cuadrado.xyz/capn-hook/core"
	"github.com/spf13/cobra"
)

var (
	hookTemplate = core.Template{
		Text: `#!/usr/bin/env bash

capn-hook run -s {hook}<<<"$(cat)"

`,
	}
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		gitDir, err := core.FindGitDir()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		for _, hookName := range core.SupportedHooks {
			hookPath := filepath.Join(gitDir, "hooks", hookName)
			fileContents := hookTemplate.Eval(core.Vars{"hook": hookName})
			ioutil.WriteFile(hookPath, []byte(fileContents), 0755)
		}

	},
}

func init() {
	RootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
