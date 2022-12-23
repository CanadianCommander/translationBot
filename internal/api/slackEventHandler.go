package api

import (
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/events"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack/slackevents"
	"io/ioutil"
	"net/http"
)

//==========================================================================
// Public
//==========================================================================

// SlackEventHandler - POST
func SlackEventHandler(gin *gin.Context) {
	log.Logger.Info("New event incoming")

	rawBody, err := ioutil.ReadAll(gin.Request.Body)
	if err != nil {
		log.Logger.Error("Failed to read request body ", err)
		gin.Status(http.StatusBadRequest)
	}

	event, err := slackevents.ParseEvent(rawBody, slackevents.OptionNoVerifyToken())
	if err != nil {
		log.Logger.Error("Failed to parse event ", err)
		gin.Status(http.StatusBadRequest)
	}

	if event.Type == slackevents.URLVerification {
		handleUrlVerification(gin, rawBody)
	} else {
		events.Dispatch(&event)
		gin.Status(http.StatusOK)
	}
}

//==========================================================================
// Private
//==========================================================================

// handleUrlVerification handle an incoming url verification challenge from slack.
// #### params
// gin - gin context
// requestBody - the raw request body
func handleUrlVerification(gin *gin.Context, requestBody []byte) {
	var challengeResponse *slackevents.ChallengeResponse
	err := json.Unmarshal([]byte(requestBody), &challengeResponse)
	if err != nil {
		gin.Status(http.StatusInternalServerError)
		return
	}

	gin.Header("Content-Type", "text")
	gin.String(http.StatusOK, challengeResponse.Challenge)
}
