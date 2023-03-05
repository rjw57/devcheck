package devcheck

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"runtime"
)

type PlatformCheck struct {
	Supported []string
}

func NewPlatformCheck(supported []string) *PlatformCheck {
	return &PlatformCheck{Supported: supported}
}

func (c *PlatformCheck) Check(l *Logger) error {
	m := mapset.NewSet[string]()
	for _, p := range c.Supported {
		m.Add(p)
	}
	if m.Contains(runtime.GOOS) {
		l.Success("%v is a platform supported by our developer environment", runtime.GOOS)
		return nil
	} else {
		l.Failure("%v is not a platform supported by our developer environment", runtime.GOOS)
		return fmt.Errorf("%v is not a supported platform", runtime.GOOS)
	}
}
