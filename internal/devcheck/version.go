package devcheck

import (
	"runtime/debug"
)

var (
	releaseRefName = ""
)

type VersionCheck struct {
	Version        string
	ReleaseRefName string
}

func NewVersionCheck() *VersionCheck {
	return &VersionCheck{}
}

func (c *VersionCheck) Check(l *Logger) error {
	c.ReleaseRefName = releaseRefName

	if c.ReleaseRefName == "" {
		l.Warning("This is not a released version of devcheck")
	} else {
		l.Info("release: %v", c.ReleaseRefName)
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				c.Version = setting.Value
			}
		}
	}

	if c.Version == "" {
		l.Warning("This version of devcheck does not have its git commit hash recorded")
	} else {
		l.Info("commit hash: %v", c.Version)
	}

	return nil
}
