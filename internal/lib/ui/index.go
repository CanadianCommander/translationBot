package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

func Index() slack.Message {
	config := configuration.Get()

	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Translation Bot :robot_face:", true, false)),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("plain_text", "How can I help you?", true, false),
			nil,
			nil),
		slack.NewDividerBlock(),
		commandOptionHelpButton("projects",
			"show a list of projects TranslationBot can translate",
			true,
			false,
			routes.ActionListProjects,
			"Run",
			""),
		commandOptionHelpButton("missing <project>",
			"show a list of missing translations",
			true,
			true,
			routes.ActionListMissingTranslations,
			"Run",
			config.DefaultProject),
		commandOptionHelpButton("all <project>",
			"get all translations as a CSV file",
			true,
			true,
			routes.ActionListAllTranslations,
			"Run",
			config.DefaultProject),
		slack.NewSectionBlock(
			slackutil.NewMarkdownTextBlock("`uploading translations` There is no command for updating translations. "+
				"Simply post a translation CSV to the translations channel and I'll see it!"),
			nil,
			nil),
		commandOptionHelpButton("notes",
			"display TranslationBot release notes",
			true,
			false,
			routes.ActionReleaseNotes,
			"Run",
			config.DefaultProject),
		slack.NewActionBlock(
			"actions",
			slack.NewButtonBlockElement(routes.ActionCancel, "", slackutil.NewTextBlock("Cancel"))),
		slack.NewContextBlock(
			"gh link",
			slackutil.NewMarkdownTextBlock("TranslationBot source code on <https://github.com/CanadianCommander/translationBot|GitHub!>")),
	)
}

// commandOptionHelp builds a new command option help message
// *command* - slash command. for example 'foobar' would be /translation foobar
// *description* - description of the command
func commandOptionHelp(command string, description string) *slack.SectionBlock {
	return commandOptionHelpButton(command, description, false, false, "", "", "")
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
	multiProject bool,
	actionId string,
	buttonTitle string,
	value string) *slack.SectionBlock {
	var accessory *slack.Accessory = nil

	if showActionButton && multiProject {
		accessory = &slack.Accessory{
			ButtonElement: slackutil.NewMultiProjectButtonBlockElement(
				actionId,
				value,
				buttonTitle,
				true),
		}
	} else if showActionButton {
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
