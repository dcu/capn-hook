// Copyright © 2016 David Cuadrado <krawek@gmail.com>
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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dcu/capn-hook/core"
	"github.com/spf13/cobra"
)

var (
	silent *bool
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <hook>",
	Short: "Runs the specified hook",
	Long: `Runs the given <hook. A hook can be either:
  ` + strings.Join(core.SupportedHooks, "\n  "),
	Run: func(cmd *cobra.Command, args []string) {
		manifest, err := core.FindManifest()
		if err != nil {
			println(err.Error())
			return
		}

		if len(args) == 0 {
			if !*silent {
				fmt.Printf("Missing hook name, options are: %v\n", core.SupportedHooks)
			}
			return
		}

		hookName := args[0]
		workingDir := filepath.Dir(manifest.Path)
		commands := manifest.Hooks(hookName)

		input := readStdin()
		for _, command := range commands {
			command.RunCommands(workingDir, input)
		}

		if len(commands) == 0 && !*silent {
			println("Invalid hook name:", hookName)
		}
	},
}

func readStdin() string {
	output := ""
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return output
	}

	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		output += reader.Text() + "\n"
	}

	return output
}

func init() {
	RootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	silent = runCmd.Flags().BoolP("silent", "s", false, "Do not print errors")
}
