package git

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"os"
	"os/exec"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

const GeneratedBranchPrefix = "translation_bot_"
const RepoStorageLocation = "./.repo/"

// InitializeProjectRepo creates & pulls the project repository and makes sure it's up-to-date & clean (no changes).
// #### params
// project - the project to pull
func InitializeProjectRepo(project *Project) error {
	checkProjectLock(project)
	log.Logger.Infof("Initializing project %s", project.Name)

	gitUrl, err := buildGitAuthUrlForProject(project)
	if err != nil {
		return err
	}

	if _, err = os.Stat(project.FilePath()); err != nil {
		log.Logger.Infof("Repository for project %s doesn't exist. Cloning", project.Name)
		cmd := exec.Command(
			"git",
			"clone",
			"-b",
			project.Branch,
			gitUrl,
			project.FilePath())
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	err = SwitchBranch(project, project.Branch, true)
	if err != nil {
		log.Logger.Errorf("Could not reset project %s to default branch %s", project.Name, project.Branch)
		return err
	}

	err = CleanProject(project)
	if err != nil {
		log.Logger.Errorf(
			"Error cleaning local state of default branch %s in project %s",
			project.Branch,
			project.Name)
		return err
	}

	err = Pull(project)
	if err != nil {
		return err
	}

	return nil
}

// Pull pulls the current branch of the specified project
// #### params
// project - the project to pull
func Pull(project *Project) error {
	log.Logger.Infof("Pulling %s", project.Name)

	cmd := exec.Command("git", "pull")
	cmd.Dir = project.FilePath()
	err := cmd.Run()
	if err != nil {
		log.Logger.Errorf("Command %s failed", cmd.String())
		return err
	}

	return nil
}

// CleanProject clears out any uncommitted changes to a project
// #### params
// project - the project to clean
func CleanProject(project *Project) error {
	log.Logger.Infof("Cleaning project %s", project.Name)

	cmd := exec.Command("git", "checkout", ".")
	cmd.Dir = project.FilePath()
	err := cmd.Run()
	if err != nil {
		log.Logger.Errorf("Command %s failed", cmd.String())
		return err
	}

	cmd = exec.Command("git", "clean", "-n", "-f")
	cmd.Dir = project.FilePath()
	err = cmd.Run()
	if err != nil {
		log.Logger.Errorf("Command %s failed", cmd.String())
		return err
	}
	return nil
}

// SwitchBranch switches the given project to the specified branch. Make sure you locked the project first!!
// #### params
// project - the project to change branches on
// branch - the branch to switch to
// pull - should the branch be pulled after switching
func SwitchBranch(project *Project, branch string, pull bool) error {
	checkProjectLock(project)
	log.Logger.Infof("Switching to branch %s", branch)

	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = project.FilePath()

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 1 {
				// branch doesn't exist. Create
				cmd = exec.Command("git", "checkout", project.Branch)
				cmd.Dir = project.FilePath()
				err = cmd.Run()
				if err != nil {
					log.Logger.Errorf("Command %s failed", cmd.String())
					return err
				}

				err = Pull(project)
				if err != nil {
					return err
				}

				cmd = exec.Command("git", "checkout", "-b", branch)
				cmd.Dir = project.FilePath()
				err = cmd.Run()
				if err != nil {
					log.Logger.Errorf("Command %s failed", cmd.String())
					return err
				}
			} else {
				log.Logger.Errorf("Command %s failed", cmd.String())
				return exitError
			}
		} else {
			log.Logger.Errorf("Command %s failed", cmd.String())
			return err
		}
	} else {
		if pull {
			return Pull(project)
		}
	}

	return nil
}

// CommitChanges commits any changes made for the specified project
func CommitChanges(project *Project) error {
	checkProjectLock(project)
	log.Logger.Infof("Commiting changes to %s", project.Name)

	// run pre commit action
	if project.PreCommitAction != "" {
		log.Logger.Infof("Running pre commit action %s", project.PreCommitAction)

		cmd := exec.Command("bash", "-c", project.PreCommitAction)
		cmd.Dir = project.FilePath()
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Logger.Errorf("Pre commit action errored out %s", err)
			log.Logger.Error(string(out))
			return err
		}
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = project.FilePath()
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", "TranslationBot. Translation update")
	cmd.Dir = project.FilePath()

	return cmd.Run()
}

// PushChanges pushes project changes to the remote
func PushChanges(project *Project) error {
	checkProjectLock(project)
	log.Logger.Infof("Pushing changes to %s", project.Name)

	currBranch, err := GetCurrentBranch(project)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "push", "--set-upstream", "origin", currBranch)
	cmd.Dir = project.FilePath()

	return cmd.Run()
}

// CommitAndPushChanges simply combines CommitChanges & CommitChanges
func CommitAndPushChanges(project *Project) error {
	checkProjectLock(project)

	err := CommitChanges(project)
	if err != nil {
		return err
	}

	return PushChanges(project)
}

// GetCurrentBranch gets the currently checked out branch for project
func GetCurrentBranch(project *Project) (string, error) {
	checkProjectLock(project)

	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = project.FilePath()

	output, err := cmd.Output()
	if err != nil {
		log.Logger.Errorf("Could not determine current branch for project %s. git commend failed", project.Name)
		return "", err
	}

	return strings.Trim(string(output), " \n"), nil
}

//==========================================================================
// Private
//==========================================================================

func checkProjectLock(project *Project) {
	if !project.IsLocked {
		panic(fmt.Sprintf("There was an attempt to modify a git project that was not locked! %s", project.Name))
	}
}
