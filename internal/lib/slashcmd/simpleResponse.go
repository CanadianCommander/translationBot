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
