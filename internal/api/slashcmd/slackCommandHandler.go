package slashcmd

import (
	"github.com/CanadianCommander/translationBot/internal/lib/slashcmd"
	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack"
	"net/http"
)

func SlashCommand(gin *gin.Context) {

	slashCommand, err := slack.SlashCommandParse(gin.Request)
	if err != nil {
		panic("Error decoding slash command")
	}

	message := slashcmd.DispatchCommand(slashCommand)
	gin.JSON(http.StatusOK, message)
}
