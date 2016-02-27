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
	"os"
	"path/filepath"

	"code.cuadrado.xyz/capn-hook/core"
	"github.com/spf13/cobra"
)

var (
	hookTemplate = `#!/usr/bin/env bash

capn-hook run -s {hook}<<<"$(cat)"

`
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs capn-hook in your git hooks",
	Long: `The install command replaces all your hooks with a capn-hook version.
If you have hooks that you are currently using please back them up before running this command.
`,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		gitDir, err := core.FindGitDir()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		for _, hookName := range core.SupportedHooks {
			hookPath := filepath.Join(gitDir, "hooks", hookName)

			os.Remove(hookPath) // In case there's a symlink

			tmpl := core.Template{Text: hookTemplate}
			tmpl.Apply(core.Vars{"hook": hookName})
			ioutil.WriteFile(hookPath, []byte(tmpl.Text), 0755)
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
