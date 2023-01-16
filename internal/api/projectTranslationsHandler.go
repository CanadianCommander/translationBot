package api

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// MissingTranslationCSVHandler - GET
func MissingTranslationCSVHandler(gin *gin.Context) {
	projectName := gin.Param("project")
	log.Logger.Infof("Incoming CSV download request for project %s", projectName)

	project := configuration.Get().GetProject(projectName)
	defer project.Unlock()

	if project == nil {
		log.Logger.Warn("Project %s not found during CSV download request", projectName)
		gin.Status(http.StatusBadRequest)
	} else {
		csvData, err := translation.FindMissingTranslationsCSV(project)

		if err != nil {
			log.Logger.Error(err)
			gin.Status(http.StatusInternalServerError)
		} else {
			responseWriter := gin.Writer
			responseWriter.Header().Add("Content-Type", "text/csv")
			responseWriter.Header().Add(
				"Content-Disposition",
				fmt.Sprintf("attachment; filename=missing_translations_%s", time.Now().Format("2006-01-02")))

			_, err = responseWriter.Write([]byte(csvData))
			if err != nil {
				log.Logger.Error("Failed to write CSV response")
			}

			responseWriter.CloseNotify()
		}
	}
}
