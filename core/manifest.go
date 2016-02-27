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
	// DefaultManifestFileName is the default name of the file that contains the manifest.
	DefaultManifestFileName = "hooks.yml"
	errManifestNotFound     = errors.New("manifest not found")
)

// Manifest represents the manifest to run the hooks
type Manifest struct {
	PreCommit        []*Hook `yaml:"pre-commit,omitempty"`
	PostReceive      []*Hook `yaml:"post-receive,omitempty"`
	PrepareCommitMsg []*Hook `yaml:"prepare-commit-msg,omitempty"`
	PostCommit       []*Hook `yaml:"post-commit,omitempty"`
	PreRebase        []*Hook `yaml:"pre-rebase,omitempty"`
	PostCheckout     []*Hook `yaml:"post-checkout,omitempty"`
	PostMerge        []*Hook `yaml:"post-merge,omitempty"`
	PrePush          []*Hook `yaml:"pre-push,omitempty"`
	PreAutoGC        []*Hook `yaml:"pre-auto-gc,omitempty"`

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

// DefaultGolangManifest returns a default manifest for golang
func DefaultGolangManifest() *Manifest {
	return &Manifest{
		PrepareCommitMsg: []*Hook{
			&Hook{
				Pattern: "*.go",
				Run: []string{
					"golint -min_confidence 0.3 {file}",
					"gocyclo -over 10 {file}",
					"varcheck",
					"deadcode",
					"structcheck",
				},
			},
		},
		PrePush: []*Hook{
			&Hook{
				Run: []string{
					"go test .",
				},
			},
		},
		PostReceive: []*Hook{
			&Hook{
				Pattern: "glide.*",
				Run: []string{
					"glide install",
				},
			},
		},
	}
}

// DefaultRubyManifest returns a default manifest for ruby
func DefaultRubyManifest() *Manifest {
	return &Manifest{
		PrepareCommitMsg: []*Hook{
			&Hook{
				Pattern: "*.rb",
				Run: []string{
					"rubycritic -f console {files}",
				},
			},
			&Hook{
				Pattern: "Gemfile*",
				Run: []string{
					"dawn -z -K .",
				},
			},
		},
		PostReceive: []*Hook{
			&Hook{
				Pattern: "Gemfile*",
				Run: []string{
					"bundle install",
				},
			},
		},
	}
}

// DefaultManifest returns a default manifest
func DefaultManifest() *Manifest {
	return &Manifest{
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*",
				Run: []string{
					"echo {files}",
					"echo {file}",
				},
			},
		},
	}
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

// ToByteArray returns the manifest encoded
func (manifest *Manifest) ToByteArray() []byte {
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return []byte{}
	}

	return data
}

// WriteFile writes the manifest to the given path
func (manifest *Manifest) WriteFile(path string) {
	ioutil.WriteFile(path, manifest.ToByteArray(), 0644)
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

		if fileName == DefaultManifestFileName {
			return LoadManifest(fileName)
		}
	}

	upDir, err := filepath.Abs(filepath.Join(path, ".."))
	if err != nil {
		return nil, err
	}

	return findManifestIn(upDir, depth+1)
}
