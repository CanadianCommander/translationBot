package api

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slashcmd"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"net/http"
)

//==========================================================================
// Public
//==========================================================================

// SlackSlashCommandHandler - POST
func SlackSlashCommandHandler(gin *gin.Context) {

	slashCommand, err := slack.SlashCommandParse(gin.Request)
	if err != nil {
		panic("Error decoding slash command")
	}

	// due to strict response time limits we spawn a new goroutine and instantly response to slack
	go slashcmd.DispatchCommand(slashCommand)
	gin.Status(http.StatusOK)
}
