package devcheck

import "os/exec"

type CommandCheck struct {
	Command    string
	InstallUrl string
	Path       string
}

func NewCommandCheck(command string) *CommandCheck {
	return &CommandCheck{Command: command}
}

func (c *CommandCheck) WithInstallUrl(url string) *CommandCheck {
	c.InstallUrl = url
	return c
}

func (c *CommandCheck) Check(l *Logger) error {
	path, err := exec.LookPath(c.Command)
	if err != nil {
		l.Failure("%v cannot be found", c.Command)
		if c.InstallUrl != "" {
			l.Indented().Info("Install instructions: %v", c.InstallUrl)
		}
		return err
	}
	c.Path = path
	l.Success("%v is installed", c.Command)
	return nil
}
