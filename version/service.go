package version

import (
	"runtime/debug"
	"sync"
)

// Version is injected during the build.
var Version = "unknown" //nolint:gochecknoglobals

type version struct {
	ApplicationVersion string `json:"applicationVersion"`
	VcsRevision        string `json:"vcsRevision"`
	VcsTime            string `json:"vcsTime"`
}

type Service interface {
	Version() version
}

func newService() Service {
	return &service{
		once: new(sync.Once),
		version: version{
			ApplicationVersion: Version,
			VcsRevision:        "unknown",
			VcsTime:            "unknown",
		},
	}
}

type service struct {
	once    *sync.Once
	version version
}

func (s *service) Version() version {
	s.once.Do(func() {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				s.version.VcsRevision = setting.Value
			} else if setting.Key == "vcs.time" {
				s.version.VcsTime = setting.Value
			}
		}
	})

	return s.version
}
