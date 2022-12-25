package actions

import (
	"errors"
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
	action := getBlockActionById(routes.ActionListMissingTranslations, interactionCallback)
	if action == nil {
		return errors.New("could not find action matching id " + routes.ActionListMissingTranslations)
	}

	bogs := translation.NewTranslation(
		"foo.bar",
		"This is english",
		[]string{"french", "japanese", "german"},
		map[string]string{"french": "la vec mwa"})

	bogs1 := translation.NewTranslation(
		"bang.boo.bazz",
		"Another one",
		[]string{"french", "japanese", "german"},
		map[string]string{"french": "la vec mwa", "japanese": "oomomi flavor"})

	response := ui.MissingTranslations([]translation.Translation{*bogs, *bogs1})

	err := slackutil.DeleteEphemeral(interactionCallback.Channel.ID, interactionCallback.ResponseURL)
	if err != nil {
		return err
	}

	_, _, err = slackutil.Api.PostMessage(interactionCallback.Channel.ID, slackutil.SlackMessageToMsgOption(response))
	if err != nil {
		return err
	}

	return nil
}
