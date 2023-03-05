package devcheck

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"os/exec"
	"regexp"
)

var MinDockerVersion = semver.MustParse("23.0.0")
var MinDockerComposeVersion = semver.MustParse("2.13.0")
var versionRegexp = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+(-[^, ]*)?`)

type DockerCheck struct {
	CommandCheck *CommandCheck
}

func NewDockerCheck(dockerCommand *CommandCheck) *DockerCheck {
	return &DockerCheck{CommandCheck: dockerCommand}
}

func (c *DockerCheck) docker(arg ...string) (string, error) {
	output, err := exec.Command(c.CommandCheck.Path, arg...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func (c *DockerCheck) versionCheck(l *Logger) error {
	versionOut, err := c.docker("--version")
	if err != nil {
		l.Failure("'docker --version' command failed: %v", err)
		return err
	}
	version, err := semver.NewVersion(versionRegexp.FindString(versionOut))
	if err != nil {
		l.Failure("Failed to parse docker version string: %v", err)
		return err
	}
	l.Success("Determined docker version: %v", version)
	if version.Compare(MinDockerVersion) >= 0 {
		l.Success("Docker version is at least %v", MinDockerVersion)
	} else {
		l.Failure("Docker version must be at least %v", MinDockerVersion)
		return fmt.Errorf("Docker version is too small")
	}
	return nil
}

func (c *DockerCheck) composeVersionCheck(l *Logger) error {
	versionOut, err := c.docker("compose", "version")
	if err != nil {
		l.Failure("'docker compose version' command failed: %v", err)
		l.Indented().Info("Docker compose plugin installtion guide: https://docs.docker.com/compose/install/")
		return err
	}
	version, err := semver.NewVersion(versionRegexp.FindString(versionOut))
	if err != nil {
		l.Failure("Failed to parse docker compose version string: %v", err)
		return err
	}
	l.Success("Determined docker compose version: %v", version)
	if version.Compare(MinDockerComposeVersion) >= 0 {
		l.Success("Docker compose version is at least %v", MinDockerComposeVersion)
	} else {
		l.Failure("Docker compose version must be at least %v", MinDockerComposeVersion)
		return fmt.Errorf("Docker compose version is too small")
	}
	return nil
}

func (c *DockerCheck) Check(l *Logger) error {
	var errs []error

	if c.CommandCheck.Path == "" {
		l.Warning("skipping since docker command was not found")
		return nil
	}

	err := c.versionCheck(l)
	if err != nil {
		errs = append(errs, err)
	}

	err = c.composeVersionCheck(l)
	if err != nil {
		errs = append(errs, err)
	}

	_, err = c.docker("run", "--rm", "hello-world")
	if err != nil {
		l.Failure("Running 'docker run --rm hello-world' failed: %v", err)
		l.Indented().Info("Make sure your user is a member of the 'docker' group")
	} else {
		l.Success("Running 'docker run --rm hello-world' succeeded")
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	} else {
		return nil
	}
}
