package slashcmd

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// simpleResponse render a simple response in response to a slash command
func simpleResponse(slashCommand *slack.SlashCommand, message slack.Message) error {
	return slackutil.PostResponse(
		slashCommand.ChannelID,
		slashCommand.ResponseURL,
		message)
}

// simplePublicResponse is just like simpleResponse but the response is public (everyone can see it)
func simplePublicResponse(slashCommand *slack.SlashCommand, message slack.Message) error {
	_, _, err := slackutil.Api.PostMessage(
		slashCommand.ChannelID,
		slackutil.SlackMessageToMsgOption(message))

	return err
}
