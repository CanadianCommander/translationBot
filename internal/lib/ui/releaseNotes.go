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

// ReleaseNotes simple page that displays shortened release notes
// #### params
// showBackButton - if true a back button will display. The button takes the user back to the index page.
func ReleaseNotes(showBackButton bool) slack.Message {
	config := configuration.Get()

	var blocks []slack.Block

	blocks = append(blocks, versionBlock(config.Version, notes2023_07_16)...)

	if showBackButton {
		blocks = append(blocks,
			slack.NewActionBlock(
				"actions",
				slack.NewButtonBlockElement(routes.ActionIndex, "", slackutil.NewTextBlock("Back")),
				slack.NewButtonBlockElement(routes.ActionReleaseNotesFull, "", slackutil.NewTextBlock("Full Notes"))),
		)
	}

	return slack.NewBlockMessage(blocks...)
}

// FullReleaseNotes simple page that displays release notes
// #### params
// showBackButton - if true a back button will display. The button takes the user back to the index page.
func FullReleaseNotes(showBackButton bool) slack.Message {
	config := configuration.Get()

	var blocks []slack.Block

	blocks = append(blocks, versionBlock("7a8cd2976b014b7117561879d10ec94a430f9572", notes2023_01_08)...)
	blocks = append(blocks, versionBlock("6a8e455191efcd7486d8bd152fdbcaceb0ee7a09", notes2023_01_15)...)
	blocks = append(blocks, versionBlock("266493d516ef4ffb6f99affe52aa9d5443df5853", notes2023_03_13)...)
	blocks = append(blocks, versionBlock(config.Version, notes2023_07_16)...)

	if showBackButton {
		blocks = append(blocks,
			slack.NewActionBlock(
				"actions",
				slack.NewButtonBlockElement(routes.ActionIndex, "", slackutil.NewTextBlock("Back"))))
	}

	return slack.NewBlockMessage(blocks...)
}

//==========================================================================
// Private
//==========================================================================

//go:embed assets/release/2023-01-08.mrkdwn
var notes2023_01_08 string

//go:embed assets/release/2023-01-15.mrkdwn
var notes2023_01_15 string

//go:embed assets/release/2023-03-13.mrkdwn
var notes2023_03_13 string

//go:embed assets/release/2023-07-16.mkrdwn
var notes2023_07_16 string

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
