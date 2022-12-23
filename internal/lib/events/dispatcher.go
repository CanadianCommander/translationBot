package events

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/slack-go/slack/slackevents"
)

//==========================================================================
// Public
//==========================================================================

// Dispatch an incoming event to the appropriate handler
// #### params
// event - incoming event
func Dispatch(event *slackevents.EventsAPIEvent) {
	log.Logger.Infof("Processing %s Event", event.InnerEvent.Type)

	switch innerEvent := event.InnerEvent.Data.(type) {
	case *slackevents.FileSharedEvent:
		OnFileShared(innerEvent)
	}
}
