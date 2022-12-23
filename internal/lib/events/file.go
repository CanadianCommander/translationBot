package events

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/slack-go/slack/slackevents"
)

//==========================================================================
// Public
//==========================================================================

// OnFileShared - respond to a new file being shared.
func OnFileShared(event *slackevents.FileSharedEvent) {
	log.Logger.Infof("User %s uploaded file %s", event.UserID, event.FileID)
}
