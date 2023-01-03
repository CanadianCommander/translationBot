package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// getBlockActionById returns the block action in the actions list that matches the provided Id or nil if there is no match
// #### params
// interactionCallback - the interaction callback containing the actions to search
func getBlockActionById(actionId string, interactionCallback *slack.InteractionCallback) *slack.BlockAction {
	for _, action := range interactionCallback.ActionCallback.BlockActions {
		if action.ActionID == actionId {
			return action
		}
	}

	return nil
}

// showLoader displays a loader to the user
// #### params
// interactionCallback - the interaction callback being processed
// msg - custom message to display on the loader
func showLoader(interactionCallback *slack.InteractionCallback, msg string) error {
	err := slackutil.PostResponse(
		interactionCallback.Channel.ID,
		interactionCallback.ResponseURL,
		ui.LoadingIndicator(msg))
	if err != nil {
		return err
	}
	return nil
}
