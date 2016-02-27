package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// PreCommitName is the name of the pre commit hook.
	PreCommitName = "pre-commit"

	// PostReceiveName is the name of the post receive hook.
	PostReceiveName = "post-receive"

	// PrepareCommitMsgName is the name of the prepare commit hook.
	PrepareCommitMsgName = "prepare-commit-msg"

	// PostCommitName is the name of the post commit hook.
	PostCommitName = "post-commit"

	// PreRebaseName is the name of the pre rebase hook.
	PreRebaseName = "pre-rebase"

	// PostCheckoutName is the name of the post checkout hook.
	PostCheckoutName = "post-checkout"

	// PostMergeName is the name of the post merge hook.
	PostMergeName = "post-merge"

	// PrePushName is the name of the pre push hook.
	PrePushName = "pre-push"

	// PreAutoGCName is the name of the pre auto GC hook.
	PreAutoGCName = "pre-auto-gc"
)

var (
	// SupportedHooks is the list of supported hooks.
	SupportedHooks = []string{
		PreCommitName, PostReceiveName, PrepareCommitMsgName, PostCheckoutName, PostCommitName, PreRebaseName, PostMergeName, PrePushName, PreAutoGCName,
	}
)

// Hook represents a hook to run
type Hook struct {
	Pattern string   `yaml:"pattern"`
	Run     []string `yaml:"run"`
}

// Match returns true if the file is matched by this hook.
func (hook *Hook) Match(filename string) (bool, error) {
	ok, err := filepath.Match(hook.Pattern, filename)
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	return filepath.Match(hook.Pattern, filepath.Base(filename))
}

// Filter filters the given files using the hook's pattern
func (hook *Hook) Filter(files []string) []string {
	filteredFiles := []string{}
	for _, file := range files {
		if file == "" {
			continue // empty entry
		}

		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue // file doesn't exist
		}

		ok, err := hook.Match(file)

		if err != nil {
			fmt.Printf("Error while matching file name %s with pattern %s: %s", file, hook.Pattern, err)
			continue
		}

		if ok {
			filteredFiles = append(filteredFiles, file)
		}

	}

	return filteredFiles
}

// RunCommand runs the given command
func (hook *Hook) RunCommand(workingDir string, command string, input string) {
	if HasTemplateVariables(command) {
		return
	}

	fmt.Printf("# Running %s\n", command)
	cmdParts := strings.Split(command, " ")

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Dir = workingDir

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("Error while running command: %s\n", err)
		os.Exit(1)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	io.Copy(stdin, bytes.NewBufferString(input))

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error while running command: %s\n", err)
		os.Exit(1)
	}
}

// RunCommands runs the command associated with this hook.
func (hook *Hook) RunCommands(workingDir string, input string) {
	for _, command := range hook.Run {
		files := FindModifiedFiles()
		filteredFiles := hook.Filter(files)

		if len(filteredFiles) == 0 && hook.Pattern != "" {
			// nothing to do here
			return
		}

		if len(filteredFiles) > 0 {
			tmpl := Template{Text: command}
			command = tmpl.Eval(Vars{"files": EscapeStringArray(filteredFiles)})
		}

		hook.RunCommand(workingDir, command, input)
	}
}

// IsSupportedHook returns true if the given hook is supported.
func IsSupportedHook(hookName string) bool {
	for _, h := range SupportedHooks {
		if hookName == h {
			return true
		}
	}

	return false
}
