package devcheck

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

const upstreamUrl = "git@github.com:rjw57/devcheck"

type FreshnessCheck struct {
	VersionCheck *VersionCheck
}

func NewFreshnessCheck(c *VersionCheck) *FreshnessCheck {
	return &FreshnessCheck{VersionCheck: c}
}

func (c *FreshnessCheck) Check(l *Logger) error {
	if c.VersionCheck.Version == "" {
		l.Warning("Skipping freshness check because devcheck version unknown")
		return nil
	}

	rem := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"git@github.com:rjw57/devcheck"},
	})

	refs, err := rem.List(&git.ListOptions{})
	if err != nil {
		l.Failure("Failed to determine the most recent devcheck version at %v: %v", upstreamUrl, err)
		return err
	}

	refsByName := make(map[string]*plumbing.Reference)
	for _, ref := range refs {
		refsByName[ref.Name().String()] = ref
	}

	headRef, headFound := refsByName[plumbing.HEAD.String()]
	for headFound && headRef.Type() == plumbing.SymbolicReference {
		headRef, headFound = refsByName[headRef.Target().String()]
	}

	if !headFound {
		l.Failure("Could not determine HEAD revision of upstream repository %v", upstreamUrl)
		return fmt.Errorf("Upstream repository %v has no HEAD", upstreamUrl)
	}

	if headRef.Hash().String() != c.VersionCheck.Version {
		l.Failure("Most recent devcheck at %v is version %v", upstreamUrl, headRef.Hash())
		return fmt.Errorf("most recent devcheck is version %v", headRef.Hash())
	}

	l.Success("Version matches most recent version from %v", upstreamUrl)
	return nil
}
