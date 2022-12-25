package ui

import "github.com/slack-go/slack"

//==========================================================================
// Public
//==========================================================================

// NotFound is a catch-all page that shows when the appropriate page is not found
func NotFound() slack.Message {

	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Translation Bot :robot_face:", true, false)),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("plain_text", "So sorry :cry: I don't understand.", true, false),
			nil,
			nil),
	)
}
