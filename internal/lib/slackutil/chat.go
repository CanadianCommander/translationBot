package slackutil

import (
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// DeleteMessageByResponseUrl removes an ephemeral slack message as identified by responseUrl
// #### params
// channelId - ID of the channel in which the ephemeral message resides
// responseUrl - the response url that identifies the ephemeral message
func DeleteMessageByResponseUrl(channelId string, responseUrl string) error {
	_, _, err := Api.PostMessage(
		channelId,
		slack.MsgOptionDeleteOriginal(responseUrl),
	)

	return err
}

// PostResponse posts a response message to slack. Replacing whatever message was previously displayed, if any.
// #### params
// channelId - the channel to post the response to
// responseUrl - slack generated url that indicates what to replace
// message - the message to post
func PostResponse(channelId string, responseUrl string, message slack.Message) error {
	_, _, err := Api.PostMessage(
		channelId,
		slack.MsgOptionReplaceOriginal(responseUrl),
		SlackMessageToMsgOption(message))
	return err
}
