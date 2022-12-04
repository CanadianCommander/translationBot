package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloWorld(g *gin.Context) {
	g.String(http.StatusOK, "Hello World")
}
