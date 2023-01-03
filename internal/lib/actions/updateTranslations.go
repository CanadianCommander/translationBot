package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

func UpdateTranslations(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) error {
	log.Logger.Infof("Applying translation update using translation file %s", block.Value)
	project := configuration.Get().GetDefaultProject()
	defer project.Unlock()

	_, err := translation.UpdateTranslationsFromSlackFile(block.Value, project)
	if err != nil {
		if err.Error() == slackutil.ErrorFileNotFound {
			err := slackutil.PostResponse(
				interactionCallback.Channel.ID,
				interactionCallback.ResponseURL,
				ui.ErrorNotification("You delete the translation file! Why would you do that?"))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}
