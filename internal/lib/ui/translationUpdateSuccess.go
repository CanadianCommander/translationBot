package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// TranslationUpdateSuccess message
// #### params
// prUrl - url of the pr in which the translation update changes can be found.
func TranslationUpdateSuccess(prUrl string) slack.Message {
	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slackutil.NewTextBlock("Translations Updated :tada:")),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slackutil.NewMarkdownTextBlock(fmt.Sprintf("Please review the changes %s", prUrl)),
			nil,
			nil),
	)
}
