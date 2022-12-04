package api

import "github.com/gin-gonic/gin"

func BuildV1Api() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("api/v1")
	v1.GET("/hello/world", HelloWorld)

	return router
}
