package actions

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ListMissingTranslations in the project specified by the block action value
func ListMissingTranslations(interactionCallback *slack.InteractionCallback) error {
	config := configuration.Get()
	action := getBlockActionById(routes.ActionListMissingTranslations, interactionCallback)
	if action == nil {
		return errors.New("could not find action matching id " + routes.ActionListMissingTranslations)
	}

	missingTranslations, err := translation.FindMissingTranslations(config.Projects[config.DefaultProject])
	if err != nil {
		return err
	}

	err = slackutil.DeleteEphemeral(interactionCallback.Channel.ID, interactionCallback.ResponseURL)
	if err != nil {
		return err
	}

	response := ui.MissingTranslations(missingTranslations)
	_, _, err = slackutil.Api.PostMessage(interactionCallback.Channel.ID, slackutil.SlackMessageToMsgOption(response))
	if err != nil {
		return err
	}

	return nil
}
