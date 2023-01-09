package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// cancelAction cancels an action workflow (deletes last message)
func cancelAction(interactionCallback *slack.InteractionCallback) error {
	return slackutil.DeleteMessageByResponseUrl(
		interactionCallback.Channel.ID,
		interactionCallback.ResponseURL)
}
