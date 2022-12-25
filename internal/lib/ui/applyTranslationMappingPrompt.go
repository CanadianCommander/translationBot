package ui

import (
	"github.com/CanadianCommander/translationBot/internal/lib/actions"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ApplyTranslationMappingsPrompt lets the user know a translation file has been detected. And asks them to confirm
// if they would like to apply the update.
func ApplyTranslationMappingsPrompt(fileId string) slack.Message {
	config := configuration.Get()

	return slack.NewBlockMessage(
		slack.NewSectionBlock(
			slackutil.NewTextBlock(":point_up: Do you want me to update the "+
				config.DefaultProject+
				" translations? "+
				slackutil.GetRandomEmoji()),
			nil,
			nil),
		slack.NewActionBlock(
			"actions",
			slack.NewButtonBlockElement(
				actions.UpdateTranslationsActionId,
				fileId,
				slackutil.NewTextBlock("Update")),
		),
	)
}
