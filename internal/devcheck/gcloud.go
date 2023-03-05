package devcheck

import (
	"fmt"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type GCloudCheck struct {
	CommandCheck *CommandCheck
}

func NewGCloudCheck(gcloudCommand *CommandCheck) *GCloudCheck {
	return &GCloudCheck{CommandCheck: gcloudCommand}
}

func (c *GCloudCheck) gcloud(arg ...string) (string, error) {
	output, err := exec.Command(c.CommandCheck.Path, arg...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

type identityResponse struct {
	Email string
}

func getIdentity(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error status calling Google API: %v", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var response identityResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Email, nil
}

func (c *GCloudCheck) Check(l *Logger) error {
	var errs []error

	if c.CommandCheck.Path == "" {
		l.Warning("skipping since gcloud command was not found")
		return nil
	}

	app_default_token, err := c.gcloud("auth", "application-default", "print-access-token")
	app_default_token = strings.TrimSpace(app_default_token)
	if err != nil {
		errs = append(errs, err)
		l.Failure("Application default credentials not configured")
		l.Indented().Info("Try using the 'gcloud auth application-default login' command")
	} else {
		l.Success("Application default credentials are configured")
		app_default_email, err := getIdentity(app_default_token)
		if err != nil {
			l.Failure("Could not determine application default credentials identity: %v", err)
			errs = append(errs, err)
		} else {
			l.Success("Application default credentials identity: %v", app_default_email)
			if !strings.HasSuffix(app_default_email, "@cam.ac.uk") {
				l.Failure("Application default credentials identity is not an @cam.ac.uk email address")
				l.Indented().Info("Try using the 'gcloud auth application-default login' command")
			}
		}
	}

	auth_token, err := c.gcloud("auth", "print-access-token")
	auth_token = strings.TrimSpace(auth_token)
	if err != nil {
		errs = append(errs, err)
		l.Failure("SDK credentials not configured")
		l.Indented().Info("Try using the 'gcloud auth login' command")
	} else {
		l.Success("SDK credentials are configured")
		auth_email, err := getIdentity(auth_token)
		if err != nil {
			l.Failure("Could not determine SDK credentials identity: %v", err)
			errs = append(errs, err)
		} else {
			l.Success("SDK credentials identity: %v", auth_email)
			if !strings.HasSuffix(auth_email, "@cam.ac.uk") {
				l.Failure("SDK credentials identity is not an @cam.ac.uk email address")
				l.Indented().Info("Try using the 'gcloud auth login' command")
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	} else {
		return nil
	}
}
