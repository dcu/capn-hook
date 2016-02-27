package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	maxDepthToFindManifest = 3
)

var (
	manifestFileName    = "hooks.yml"
	errManifestNotFound = errors.New("Manifest not found")
)

// Manifest represents the manifest to run the hooks
type Manifest struct {
	PreCommit        []*Hook `yaml:"pre-commit"`
	PostReceive      []*Hook `yaml:"post-receive"`
	PrepareCommitMsg []*Hook `yaml:"prepare-commit-msg"`
	PostCommit       []*Hook `yaml:"post-commit"`
	PreRebase        []*Hook `yaml:"pre-rebase"`
	PostCheckout     []*Hook `yaml:"post-checkout"`
	PostMerge        []*Hook `yaml:"post-merge"`
	PrePush          []*Hook `yaml:"pre-push"`
	PreAutoGC        []*Hook `yaml:"pre-auto-gc"`

	Path string `yaml:"-"`
}

// LoadManifest loads the manifest from the given path
func LoadManifest(path string) (*Manifest, error) {
	manifest := &Manifest{Path: path}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

// FindManifest finds the manifest by navigating the parent directories.
func FindManifest() (*Manifest, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return findManifestIn(wd, 0)
}

// Hooks returns all associated hooks given a hook name.
func (manifest *Manifest) Hooks(name string) []*Hook {
	switch name {
	case PreCommitName:
		{
			return manifest.PreCommit
		}
	case PostReceiveName:
		{
			return manifest.PostReceive
		}
	case PrepareCommitMsgName:
		{
			return manifest.PrepareCommitMsg
		}
	case PostCommitName:
		{
			return manifest.PostCommit
		}
	case PreRebaseName:
		{
			return manifest.PreRebase
		}
	case PostCheckoutName:
		{
			return manifest.PostCheckout
		}
	case PostMergeName:
		{
			return manifest.PostMerge
		}
	case PrePushName:
		{
			return manifest.PrePush
		}
	case PreAutoGCName:
		{
			return manifest.PreAutoGC
		}
	}

	return nil
}

func findManifestIn(path string, depth int) (*Manifest, error) {
	if depth > maxDepthToFindManifest {
		return nil, errManifestNotFound
	}

	matches, err := filepath.Glob(filepath.Join(path, "*.yml"))
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		fileName := filepath.Base(match)

		if fileName == manifestFileName {
			return LoadManifest(fileName)
		}
	}

	upDir, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return nil, err
	}

	return findManifestIn(upDir, depth+1)
}
