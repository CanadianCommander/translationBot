package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

func Index() slack.Message {

	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Translation Bot :robot_face:", true, false)),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("plain_text", "How can I help you?", true, false),
			nil,
			nil),
		slack.NewDividerBlock(),
		commandOptionHelpButton("missing <en|fr>", "show a list of missing translations", true, "TODO", "Run", "fr"),
		commandOptionHelpButton("upload", "Upload a new translation file (CSV). \n"+
			"I can also automatically detect translations files uploaded to slack :wink: Try uploading one in this channel.",
			true,
			"TODO",
			"Run",
			""),
	)
}

// commandOptionHelp builds a new command option help message
// *command* - slash command. for example 'foobar' would be /translation foobar
// *description* - description of the command
func commandOptionHelp(command string, description string) *slack.SectionBlock {
	return commandOptionHelpButton(command, description, false, "", "", "")
}

// commandOptionHelpButton builds a new command option help message
// *command* - slash command. for example 'foobar' would be /translation foobar
// *description* - description of the command
// *showActionButton* - should an action button be shown.
// *actionId* - action id to invoke if the user clicks the action button.
// *buttonTitle* - the title of the button
// *value* - value to send to the action
func commandOptionHelpButton(
	command string,
	description string,
	showActionButton bool,
	actionId string,
	buttonTitle string,
	value string) *slack.SectionBlock {
	var accessory *slack.Accessory = nil

	if showActionButton {
		accessory = &slack.Accessory{
			ButtonElement: slack.NewButtonBlockElement(
				actionId,
				value,
				slackutil.NewTextBlock(buttonTitle)),
		}
	}

	return slack.NewSectionBlock(
		slackutil.NewMarkdownTextBlock(fmt.Sprintf("`/translation %s` %s", command, description)),
		nil,
		accessory,
	)
}
