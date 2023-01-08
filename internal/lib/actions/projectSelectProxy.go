package actions

import (
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
	"golang.org/x/exp/maps"
)

//==========================================================================
// Private
//==========================================================================

// projectSelectProxy is an action that allows the user to select the desired project. The original value passed to this action
// is then passed along to the real action wrapped in a struct indicating which project was selected
func projectSelectProxy(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) error {
	config := configuration.Get()
	projects := maps.Values(config.Projects)

	input := slackutil.MultiProjectInputDto{}
	err := json.Unmarshal([]byte(block.Value), &input)
	if err != nil {
		return err
	}

	actions := make([]string, len(projects))
	values := make([]string, len(projects))
	for idx, project := range projects {
		actions[idx] = input.ActionId

		selectProxyReq := ProjectSelectProxyRequestDto{
			Project: project.Name,
			Value:   input.Value,
		}

		jsonBytes, err := json.Marshal(selectProxyReq)
		if err != nil {
			return err
		}

		values[idx] = string(jsonBytes)
	}

	return slackutil.PostResponse(
		interactionCallback.Channel.ID,
		interactionCallback.ResponseURL,
		ui.ProjectSelectList(projects, actions, values))
}
