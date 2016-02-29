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
			},
		},
	}
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
