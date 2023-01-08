package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
	"golang.org/x/exp/maps"
)

//==========================================================================
// Private
//==========================================================================

// listProjects list projects
func listProjects(interactionCallback *slack.InteractionCallback) error {
	config := configuration.Get()

	err := slackutil.PostResponse(
		interactionCallback.Channel.ID,
		interactionCallback.ResponseURL,
		ui.ProjectList(maps.Values(config.Projects)))
	return err
}
