package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// Dispatch incoming action to appropriate handler
func Dispatch(interactionCallback *slack.InteractionCallback) error {

	for _, block := range interactionCallback.ActionCallback.BlockActions {
		go dispatchAction(interactionCallback, block)
	}

	return nil
}

//==========================================================================
// Private
//==========================================================================

// dispatchAction dispatches an action to the appropriate handler
func dispatchAction(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) {
	var err error = nil

	log.Logger.Infof("Handling action %s", block.ActionID)
	switch block.ActionID {
	case routes.ActionListMissingTranslations:
		err = ListMissingTranslations(interactionCallback)
	case routes.ActionUpdateTranslations:
		err = UpdateTranslations(interactionCallback)
	}

	if err != nil {
		log.Logger.Error("Unexpected error while handling slack action ", err)
	}
}
