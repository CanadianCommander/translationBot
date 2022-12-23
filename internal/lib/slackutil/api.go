package slackutil

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

var Api *slack.Client = nil

// InitializeApi setups the slack api connection
func InitializeApi() {
	Api = slack.New(configuration.Get().SlackClientSecret)
}
