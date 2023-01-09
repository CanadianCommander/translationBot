package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// simpleAction renders a fixed response
func simpleAction(interactionCallback *slack.InteractionCallback, msg slack.Message) error {

	return slackutil.PostResponse(
		interactionCallback.Channel.ID,
		interactionCallback.ResponseURL,
		msg)
}

// simpleActionNoReplace is like simple action, but it does not replace the existing slack message
func simpleActionNoReplace(interactionCallback *slack.InteractionCallback, msg slack.Message) error {
	_, _, err := slackutil.Api.PostMessage(
		interactionCallback.Channel.ID,
		slackutil.SlackMessageToMsgOption(msg))

	return err
}
