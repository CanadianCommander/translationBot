package gh

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/google/go-github/v48/github"
	"time"
)

//==========================================================================
// Public
//==========================================================================

// CreatePr creates a pr in the given project for the given branch.
// #### params
// project - the project to create the PR in
// branch - the branch to merge in to the projects default branch (as configured in config file)
// #### return
// url where the pull request can be viewed
func CreatePr(project *git.Project, branch string) (string, error) {
	log.Logger.Infof("Creating PR in project %s for branch %s", project.Name, branch)

	client, ghAuthCtx := GetClientForProject(project)
	owner, repo, err := getGitHubOwnerRepoFromUrl(project.Url)
	if err != nil {
		return "", err
	}

	title := fmt.Sprintf("Translation Update %s", time.Now().Format("2006-01-02"))
	body := "## :robot: Another translation update from your friendly neighborhood TranslationBot\n" +
		"**Beep Boop**"

	pr, _, err := client.PullRequests.Create(
		ghAuthCtx,
		owner,
		repo,
		&github.NewPullRequest{
			Title: &title,
			Head:  &branch,
			Base:  &project.Branch,
			Body:  &body,
		},
	)
	if err != nil {
		return "", err
	}

	return *pr.HTMLURL, nil
}
