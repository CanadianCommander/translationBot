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
		var err error = nil

		log.Logger.Infof("Handling action %s", block.ActionID)
		switch block.ActionID {
		case routes.ActionListMissingTranslations:
			err = ListMissingTranslations(interactionCallback)
		case routes.ActionUpdateTranslations:
			err = UpdateTranslations(interactionCallback)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
