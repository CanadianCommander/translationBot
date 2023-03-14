package api

import (
	"github.com/CanadianCommander/translationBot/internal/api/middleware/security"
	"github.com/gin-gonic/gin"
)

func BuildV1Api() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("api/v1/")

	v1Slack := v1.Group("slack/")
	v1Slack.Use(security.ValidateSlackRequest)
	v1Slack.POST("/cmd/", SlackSlashCommandHandler)
	v1Slack.POST("/event/", SlackEventHandler)
	v1Slack.POST("/action/", SlackActionHandler)

	// downloads
	v1Project := v1.Group("/project/:project/")
	v1Project.GET("/translations/missing/csv", MissingTranslationCSVHandler)
	v1Project.GET("/translations/csv", TranslationCSVHandler)

	return router
}
