package devcheck

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
)

const (
	repoOwner = "rjw57"
	repoName  = "devcheck"
)

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

	ctx := context.Background()
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, repoOwner, repoName)
	if err != nil {
		l.Failure("Could not get latest release from https://github.com/%v/%v: %v",
			repoOwner, repoName, err)
		return err
	}

	ref, _, err := client.Git.GetRef(ctx, repoOwner, repoName, "refs/tags/"+*release.TagName)
	if err != nil {
		l.Failure(
			"Could not get SHA for latest release '%v' from https://github.com/%v/%v: %v",
			*release.TagName, repoOwner, repoName, err)
		return err
	}

	if *ref.Object.SHA != c.VersionCheck.Version {
		l.Failure(
			"Most recent release of devcheck at https://github.com/%v/%v is '%v'",
			repoOwner, repoName, *release.TagName)
		l.Indented().Info("Download the most recent release at: %v", *release.HTMLURL)
		return fmt.Errorf("most recent devcheck is version %v", *release.TagName)
	}

	l.Success(
		"Version matches most recent version '%v' from https://github.com/%v/%v",
		*release.TagName, repoOwner, repoName)
	return nil
}
