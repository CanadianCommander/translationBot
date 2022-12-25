package api

import (
	"github.com/CanadianCommander/translationBot/internal/api/middleware/security"
	"github.com/gin-gonic/gin"
)

func BuildV1Api() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("api/v1")

	v1.Use(security.ValidateSlackRequest)
	v1.POST("/cmd/", SlackSlashCommandHandler)
	v1.POST("/event/", SlackEventHandler)
	v1.POST("/action/", SlackActionHandler)

	return router
}
