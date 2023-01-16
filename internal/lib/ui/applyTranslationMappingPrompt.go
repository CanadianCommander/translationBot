package ui

import (
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ApplyTranslationMappingsPrompt lets the user know a translation file has been detected. And asks them to confirm
// if they would like to apply the update.
func ApplyTranslationMappingsPrompt(fileId string) slack.Message {
	return slack.NewBlockMessage(
		slack.NewSectionBlock(
			slackutil.NewTextBlock(":point_up: Do you want me to update the translations? "+
				slackutil.GetRandomEmoji()),
			nil,
			nil),
		slack.NewActionBlock(
			"actions",
			slack.NewButtonBlockElement(
				routes.ActionCancel,
				"",
				slackutil.NewTextBlock("Dismiss")),
			slackutil.NewMultiProjectButtonBlockElement(
				routes.ActionUpdateTranslations,
				fileId,
				"Update",
				false),
		),
	)
}
