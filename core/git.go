package core

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	maxDepthToFindGitDir = 2
)

var (
	errGitDirNotFound = errors.New("GITDIR not found")
	gitDirFiles       = []string{"HEAD", "config", "description", "index", "objects", "hooks"}
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

// FindGitDir finds the path to the GITDIR
func FindGitDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return findGitDirIn(wd, 0)
}

func findGitDirIn(path string, depth int) (string, error) {
	if depth > maxDepthToFindGitDir {
		return "", errGitDirNotFound
	}

	candidate := filepath.Join(path, ".git")
	isGitDir := true

	for _, gitDirFile := range gitDirFiles {
		file := filepath.Join(candidate, gitDirFile)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			isGitDir = false
			break
		}
	}

	if isGitDir {
		return candidate, nil
	}

	upDir, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return "", err
	}

	return findGitDirIn(upDir, depth+1)
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
