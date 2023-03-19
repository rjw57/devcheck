package devcheck

import (
	"context"

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
	ctx := context.Background()
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(ctx, repoOwner, repoName)
	if err != nil {
		l.Failure("Could not get latest release from https://github.com/%v/%v: %v",
			repoOwner, repoName, err)
		return err
	}

	if c.VersionCheck.ReleaseRefName == "" {
		l.Warning("Skipping freshness check because this is not a released version")
	} else {
		if *release.TagName != c.VersionCheck.ReleaseRefName {
			l.Warning("The most recent release is '%v'", *release.TagName)
			l.Indented().Info(
				"Download the most recent release at: %v", *release.HTMLURL)
		} else {
			l.Success("This is the most recent release")
		}

		if c.VersionCheck.Version == "" {
			l.Warning("Skipping version check because devcheck version unknown")
		} else {
			ref, _, err := client.Git.GetRef(
				ctx, repoOwner, repoName, "refs/tags/"+c.VersionCheck.ReleaseRefName)
			if err != nil {
				l.Failure(
					"Could not get SHA for release '%v' from https://github.com/%v/%v: %v",
					c.VersionCheck.ReleaseRefName, repoOwner, repoName, err)
				return err
			}

			if *ref.Object.SHA != c.VersionCheck.Version {
				l.Failure("Commit hash for release '%v' should be %v",
					c.VersionCheck.ReleaseRefName, *ref.Object.SHA)
			} else {
				l.Success("Commit hash matches release tag")
			}
		}
	}
	return nil
}
