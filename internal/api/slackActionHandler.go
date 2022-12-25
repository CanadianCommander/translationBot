package api

import (
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/actions"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"net/http"
)

//==========================================================================
// Public
//==========================================================================

// SlackActionHandler - POST
func SlackActionHandler(gin *gin.Context) {
	interactionCallback, err := parseInteractionCallbackPayload(gin)

	if err != nil {
		log.Logger.Error("Error parsing slack action event ", err)
		gin.Status(http.StatusBadRequest)
	} else {
		log.Logger.Infof("Incomming action %s", interactionCallback.Type)

		if err := actions.Dispatch(interactionCallback); err != nil {
			log.Logger.Error("Unexpected error while processing block action ", err)
			gin.Status(http.StatusInternalServerError)
		} else {
			gin.Status(http.StatusOK)
		}
	}
}

//==========================================================================
// Private
//==========================================================================

// parseInteractionCallbackPayload provides functionality similar to slack.parseSlashCommand and
// slackevent.ParseEvent. There is an issue open to add this feature to slack-go
// https://github.com/slack-go/slack/issues/660
func parseInteractionCallbackPayload(gin *gin.Context) (*slack.InteractionCallback, error) {
	interactionCallback := slack.InteractionCallback{}
	err := json.Unmarshal([]byte(gin.Request.FormValue("payload")), &interactionCallback)
	if err != nil {
		return nil, err
	}
	return &interactionCallback, nil
}
