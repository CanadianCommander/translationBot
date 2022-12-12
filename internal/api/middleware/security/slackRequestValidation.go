package security

import (
	"bytes"
	"github.com/CanadianCommander/translationBot/internal/lib/api"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"io/ioutil"
	"net/http"
)

// ValidateSlackRequest
// Validates that an incoming request is really from slack.
// This method panics if the request did not come from slack servers.
func ValidateSlackRequest(gin *gin.Context) {
	config := configuration.Load()

	sv, secretErr := slack.NewSecretsVerifier(gin.Request.Header, config.SlackSigningSecret)
	body, readErr := ioutil.ReadAll(gin.Request.Body)
	if readErr != nil || secretErr != nil {
		log.Logger.Warn("Failed to read request body or build validator")
		gin.JSON(http.StatusUnauthorized, api.NewUnauthorizedError())
		gin.Abort()
	}

	_, secretWriteError := sv.Write(body)
	if secretWriteError != nil || sv.Ensure() != nil {
		gin.JSON(http.StatusUnauthorized, api.NewUnauthorizedError())
		gin.Abort()
	}

	// Write body back to Gin. Gin streams data in only once from the client. Because we need to verify the body
	// we must buffer the data. Sort of a hack to deal with the limited compatability between slack-go and gin
	gin.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}
