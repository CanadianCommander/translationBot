package slackutil

import "github.com/slack-go/slack"

//==========================================================================
// Public
//==========================================================================

// DeleteEphemeral removes an ephemeral slack message as identified by responseUrl
// #### params
// channelId - ID of the channel in which the ephemeral message resides
// responseUrl - the response url that identifies the ephemeral message
func DeleteEphemeral(channelId string, responseUrl string) error {
	_, _, err := Api.PostMessage(
		channelId,
		slack.MsgOptionDeleteOriginal(responseUrl),
	)

	return err
}
