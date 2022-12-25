package actions

import "github.com/slack-go/slack"

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
