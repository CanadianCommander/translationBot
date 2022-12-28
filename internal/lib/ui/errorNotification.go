package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ErrorNotification shows a nice little error message to the user.
func ErrorNotification(msg string) slack.Message {
	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slackutil.NewTextBlock("That's not right!?!?!")),
		slack.NewSectionBlock(
			slackutil.NewTextBlock(fmt.Sprintf("I'm sorry I can't complete your request right now.\n %s", msg)),
			nil,
			nil))
}
