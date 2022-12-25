package events

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translationMapping"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack/slackevents"
)

//==========================================================================
// Public
//==========================================================================

// OnFileShared - respond to a new file being shared.
func OnFileShared(event *slackevents.FileSharedEvent) {
	log.Logger.Infof("User %s uploaded file %s checking type", event.UserID, event.FileID)

	loader, err := translationMapping.GetMapperForSlackFile(event.FileID)
	if err != nil {
		log.Logger.Error("Unexpected error when searching for file loader")
	}

	if loader != nil {
		log.Logger.Infof("File looks like translation mapping file! Matches loader %s", loader.GetLoaderType())
		_, _, err := slackutil.Api.PostMessage(
			event.ChannelID,
			slackutil.SlackMessageToMsgOption(ui.ApplyTranslationMappingsPrompt(event.FileID)))
		if err != nil {
			log.Logger.Error("Unexpected error posting translation mapping prompt to slack ", err)
		}
	} else {
		log.Logger.Info("File doesn't look like translation mapping file. Ignoring")
	}
}
