package devcheck

import (
	"fmt"
	"runtime/debug"
)

type VersionCheck struct {
	Version string
}

func NewVersionCheck() *VersionCheck {
	return &VersionCheck{Version: version()}
}

func (c *VersionCheck) Check(l *Logger) error {
	if c.Version == "" {
		l.Failure("Could not determine devcheck version")
		return fmt.Errorf("Could not determine devcheck version")
	}

	l.Success("Determined devcheck version: %v", c.Version)
	return nil
}

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}

	return ""
}
