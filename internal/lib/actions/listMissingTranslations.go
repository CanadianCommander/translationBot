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
func ListMissingTranslations(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) error {
	projectName, _, err := readProxyResponse(block.Value)
	if err != nil {
		return err
	}

	project := configuration.Get().GetProject(projectName)
	defer project.Unlock()

	action := getBlockActionById(routes.ActionListMissingTranslations, interactionCallback)
	if action == nil {
		return errors.New("could not find action matching id " + routes.ActionListMissingTranslations)
	}

	err = showLoader(interactionCallback, "Degaussing the translation matrix...")
	if err != nil {
		return err
	}

	missingTranslations, err := translation.FindMissingTranslations(project)
	if err != nil {
		return err
	}

	err = slackutil.DeleteMessageByResponseUrl(interactionCallback.Channel.ID, interactionCallback.ResponseURL)
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
