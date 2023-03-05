package devcheck

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"regexp"
)

type SSHConnection struct {
	User string
	Addr string
}

type SSHCheck struct {
	Connections []string
}

func NewSSHCheck(arg ...string) *SSHCheck {
	return &SSHCheck{Connections: arg}
}

func (s *SSHCheck) Add(connection string) *SSHCheck {
	s.Connections = append(s.Connections, connection)
	return s
}

func (s *SSHCheck) Check(l *Logger) error {
	socket := os.Getenv("SSH_AUTH_SOCK")

	conn, err := net.Dial("unix", socket)
	if err != nil {
		l.Failure("Could not connect to SSH agent")
		return err
	}
	agentClient := agent.NewClient(conn)

	signers, err := agentClient.Signers()
	if err != nil {
		l.Failure("Could not fetch private keys from SSH agent")
		return err
	}

	if len(signers) > 0 {
		l.Success("At least one key available in SSH agent")
	} else {
		l.Failure("No keys in SSH agent")
		return fmt.Errorf("No keys in SSH agent")
	}

	var errs []error
	connectionRegexp := regexp.MustCompile(
		`(?P<user>[a-zA-Z0-9-\.]+)@(?P<host>[a-zA-Z0-9-,\.]+)`)

	for _, c := range s.Connections {
		submatches := connectionRegexp.FindStringSubmatch(c)
		if submatches == nil {
			l.Failure("Failed to parse connection specifier: %v", c)
			continue
		}
		user := submatches[connectionRegexp.SubexpIndex("user")]
		host := submatches[connectionRegexp.SubexpIndex("host")]
		config := &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signers...)},
			// TODO: parse known_hosts file if present
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		sshc, err := ssh.Dial("tcp", host+":22", config)
		if err != nil {
			l.Failure("%v: login failed: %v", c, err)
			l.Indented().Info("Check that your key has been added to the SSH agent via ssh-add")
			errs = append(errs, err)
			break
		}
		l.Success("%v: login succeeded", c)
		sshc.Close()
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	} else {
		return nil
	}
}
