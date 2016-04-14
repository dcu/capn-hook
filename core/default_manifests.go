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
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*.go",
				Run: []string{
					"golint -min_confidence 0.3 {file}",
					"gocyclo -over 10 {file}",
					"varcheck",
					"deadcode",
					"structcheck",
				},
				Enforce: true,
			},
		},
		PrePush: []*Hook{
			&Hook{
				Run: []string{
					"go test .",
				},
				Enforce: true,
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
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*.rb",
				Run: []string{
					"rubycritic -f console {files}",
				},
				Enforce: true,
			},
			&Hook{
				Pattern: "Gemfile*",
				Run: []string{
					"dawn -z -K .",
				},
				Enforce: true,
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

// DefaultAndroidManifest returns the default manifest for android
func DefaultAndroidManifest() *Manifest {
	return &Manifest{
		PreCommit: []*Hook{
			&Hook{
				Pattern: "*.java",
				Run: []string{
					"lint .",
				},
			},
			&Hook{
				Pattern: "*.xml",
				Run: []string{
					"lint .",
				},
			},
		},
	}
}
