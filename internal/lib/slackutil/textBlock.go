package slackutil

import "github.com/slack-go/slack"

// NewTextBlock makes a new slack text block
func NewTextBlock(text string) *slack.TextBlockObject {
	return slack.NewTextBlockObject("plain_text", text, true, false)
}

// NewMarkdownTextBlock makes a new markdown slack text block
func NewMarkdownTextBlock(text string) *slack.TextBlockObject {
	return slack.NewTextBlockObject("mrkdwn", text, false, false)
}
