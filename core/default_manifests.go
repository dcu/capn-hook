package core

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
				Required: false,
			},
		},
	}
}

// DefaultGolangManifest returns a default manifest for golang
func DefaultGolangManifest() *Manifest {
	return &Manifest{
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*.go",
				Run: []string{
					"golint -min_confidence 0.3 -set_exit_status {files}",
					"gocyclo -over 10 {file}",
					"varcheck",
					"deadcode",
					"structcheck",
				},
				Required: true,
			},
		},
		PrePush: []*Hook{
			&Hook{
				Run: []string{
					"go test .",
				},
				Required: true,
			},
		},
		PostReceive: []*Hook{
			&Hook{
				Pattern: "glide.*",
				Run: []string{
					"glide install",
				},
				Required: false,
			},
		},
	}
}

// DefaultRubyManifest returns a default manifest for ruby
func DefaultRubyManifest() *Manifest {
	return &Manifest{
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*.rb",
				Run: []string{
					"rubycritic -f console {files}",
				},
				Required: true,
			},
			&Hook{
				Pattern: "Gemfile*",
				Run: []string{
					"dawn -z -K .",
				},
				Required: true,
			},
		},
		PostReceive: []*Hook{
			&Hook{
				Pattern: "Gemfile*",
				Run: []string{
					"bundle install",
				},
				Required: false,
			},
		},
	}
}

// DefaultAndroidManifest returns the default manifest for android
func DefaultAndroidManifest() *Manifest {
	return &Manifest{
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*.java",
				Run: []string{
					"lint .",
				},
				Required: false,
			},
			&Hook{
				Pattern: "*.xml",
				Run: []string{
					"lint .",
				},
				Required: false,
			},
		},
	}
}
