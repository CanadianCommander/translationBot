package git

import (
	"errors"
	"fmt"
	"strings"
)

//==========================================================================
// Private
//==========================================================================

// buildGitAuthUrlForProject creates an authenticated git url for the given project.
// this allows for pulling of private repositories. If auth details are not included in project
// then the original git url is returned.
func buildGitAuthUrlForProject(project *Project) (string, error) {
	if project.GitUsername != "" && project.GitPassword != "" {
		urlProtoSplit := strings.Split(project.Url, "://")
		if len(urlProtoSplit) != 2 {
			return "", errors.New("cannot build auth url, git url malformed: " + project.Url)
		}

		return fmt.Sprintf("%s://%s:%s@%s", urlProtoSplit[0], project.GitUsername, project.GitPassword, urlProtoSplit[1]), nil
	}
	return project.Url, nil
}
