package git

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"os"
	"os/exec"
)

//==========================================================================
// Public
//==========================================================================

const RepoStorageLocation = "/tmp/"

// PullProjectRepo pulls the project repository and makes sure it's up-to-date.
// #### params
// project - the project to pull
func PullProjectRepo(project *Project) error {
	gitUrl, err := buildGitAuthUrlForProject(project)
	if err != nil {
		return err
	}

	if _, err := os.Stat(project.filePath()); err != nil {
		log.Logger.Infof("Repository for project %s doesn't exist. Cloning", project.Name)
		err = exec.Command(
			"git",
			"clone",
			"-b",
			project.Branch,
			gitUrl,
			project.filePath()).Run()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("git", "pull", gitUrl)
	cmd.Dir = project.filePath()

	log.Logger.Infof("Pulling down project updates for %s", project.Name)
	return cmd.Run()
}
