package main

import dc "github.com/rjw57/devcheck/internal/devcheck"

func main() {
	dockerCmdCheck := dc.NewCommandCheck("docker").WithInstallUrl("https://docs.docker.com/get-docker/")
	gcloudCmdCheck := dc.NewCommandCheck("gcloud").WithInstallUrl("https://cloud.google.com/sdk/docs/install")

	dc.NewCheckList(
		dc.NewSectionCheck(
			"Checking platform",
			dc.NewPlatformCheck([]string{"linux", "darwin"}),
		),
		dc.NewSectionCheck("Checking installed software", dc.NewCheckList(
			dockerCmdCheck,
			gcloudCmdCheck,
			dc.NewCommandCheck("python3").WithInstallUrl("https://wiki.python.org/moin/BeginnersGuide/Download"),
			dc.NewCommandCheck("git").WithInstallUrl("https://wiki.python.org/moin/BeginnersGuide/Download"),
			dc.NewCommandCheck("ssh"),
		)),
		dc.NewSectionCheck("Checking Docker", dc.NewDockerCheck(dockerCmdCheck)),
		dc.NewSectionCheck("Checking SSH", dc.NewSSHCheck(
			"git@github.com",
			"git@gitlab.developers.cam.ac.uk",
		)),
		dc.NewSectionCheck("Checking Google Cloud SDK", dc.NewGCloudCheck(gcloudCmdCheck)),
	).Check(dc.NewLogger())
}
