package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// CsvDownload produces a simple UI allowing the user to initiate a CSV download
// #### params
// projectName - the name of the project from which the CSV download comes
// downloadPath - a relative download path. The part after, /api/v1/project/<project>/
func CsvDownload(projectName string, downloadPath string) slack.Message {
	config := configuration.Get()

	return slack.NewBlockMessage(
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slackutil.NewTextBlock("Click here to download CSV"),
			nil,
			slack.NewAccessory(
				&slack.ButtonBlockElement{
					Type: "button",
					Text: slackutil.NewTextBlock("CSV"),
					URL:  fmt.Sprintf("https://%s/api/v1/project/%s/%s", config.Hostname, projectName, downloadPath),
				}),
		))
}
