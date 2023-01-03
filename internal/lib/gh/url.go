package gh

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

//==========================================================================
// Url
//==========================================================================

// gitGitHubOwnerRepoFromUrl splits the owner and the repo out of a GitHub url.
// #### params
// githubUrl - the url to split
// #### return
// owner, repo. Error will be set if the url is not a github url
func getGitHubOwnerRepoFromUrl(githubUrl string) (string, string, error) {
	if !strings.Contains(githubUrl, "github.com") {
		return "", "", errors.New(fmt.Sprintf("cannot extract owner/repo! %s is not a github url", githubUrl))
	}

	expression, err := regexp.Compile("([^/]+)/([^./]+)[^/]+$")
	if err != nil {
		return "", "", err
	}

	matches := expression.FindStringSubmatch(githubUrl)
	if len(matches) != 3 {
		return "", "", errors.New(fmt.Sprintf("got unexpected number of path segements when splitting git url %s", githubUrl))
	}

	return matches[1], matches[2], nil
}
