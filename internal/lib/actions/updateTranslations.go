package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/gh"
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

	config := configuration.Get()
	project := config.GetDefaultProject()
	defer project.Unlock()

	err := showLoader(interactionCallback, "Loading Universal Translator...")
	if err != nil {
		return err
	}

	_, err = translation.UpdateTranslationsFromSlackFile(block.Value, project)
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

	prUrl := "N/A"
	if !config.TestMode {
		prUrl, err = gh.CreatePr(project, "translation_bot_2023-01-03_1672706791")
		if err != nil {
			return err
		}
	}

	err = slackutil.PostResponse(
		interactionCallback.Channel.ID,
		interactionCallback.ResponseURL,
		ui.TranslationUpdateSuccess(prUrl))
	if err != nil {
		return err
	}

	return nil
}
