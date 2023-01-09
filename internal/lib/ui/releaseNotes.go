package ui

import (
	_ "embed"
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ReleaseNotes simple page that displays release notes
func ReleaseNotes() slack.Message {
	config := configuration.Get()

	var blocks []slack.Block

	blocks = append(blocks, versionBlock(config.Version, notes2023_01_08)...)
	blocks = append(blocks,
		slack.NewActionBlock(
			"actions",
			slack.NewButtonBlockElement(routes.ActionIndex, "", slackutil.NewTextBlock("Back"))))

	return slack.NewBlockMessage(blocks...)
}

//==========================================================================
// Private
//==========================================================================

//go:embed assets/release/2023-01-08.mrkdwn
var notes2023_01_08 string

// versionBlock renders blocks to display release notes for a version of TranslationBot
func versionBlock(version string, notes string) []slack.Block {

	return []slack.Block{
		slack.NewHeaderBlock(slackutil.NewTextBlock(fmt.Sprintf("Version %s", version))),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slackutil.NewMarkdownTextBlock(notes),
			nil,
			nil),
	}
}
