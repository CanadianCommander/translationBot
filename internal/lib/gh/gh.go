package gh

import (
	"context"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

//==========================================================================
// Public
//==========================================================================

// GetClientForProject returns a github client that can be used to interact with the given project
// #### params
// project - the project to get a github client for.
func GetClientForProject(project *git.Project) (*github.Client, context.Context) {

	ctx := context.Background()
	staticToken := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: project.GitPassword})

	return github.NewClient(oauth2.NewClient(ctx, staticToken)), ctx
}
