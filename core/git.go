package core

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

// GitCommand is a command to be executed by git
type GitCommand struct {
	ProcInput *bytes.Reader
	Args      []string
}

// Run runs the git command
func (gitCommand *GitCommand) Run(wait bool) (io.ReadCloser, error) {
	cmd := exec.Command("git", gitCommand.Args...)
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	if gitCommand.ProcInput != nil {
		cmd.Stdin = gitCommand.ProcInput
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if wait {
		err = cmd.Wait()
		if err != nil {
			return nil, err
		}
	}

	return stdout, nil
}

// RunAndGetOutput runs the command and gets the output
func (gitCommand *GitCommand) RunAndGetOutput() []byte {
	stdout, err := gitCommand.Run(false)
	if err != nil {
		return []byte{}
	}

	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return []byte{}
	}

	return data
}

// FindModifiedFiles returns the list of all modified files
func FindModifiedFiles() []string {
	result := GitDiff("--name-only", "-z")
	result = append(result, GitDiff("--name-only", "--cached", "-z")...)

	return result
}

// GitDiff runs the git-diff command
func GitDiff(options ...string) []string {
	command := &GitCommand{Args: append([]string{"diff"}, options...)}

	output := command.RunAndGetOutput()
	return strings.Split(string(output), "\x00")
}
