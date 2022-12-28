package ui

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

// LoadingIndicator render a loading indicator
// #### params
// subheading - sub text to be displayed bellow the pleasing wait message.
func LoadingIndicator(subheading string) slack.Message {
	return slack.NewBlockMessage(
		slack.NewSectionBlock(
			slackutil.NewTextBlock("Please wait... :loading:"),
			nil,
			nil),
		slack.NewContextBlock(
			"",
			slackutil.NewTextBlock(subheading)),
	)
}
