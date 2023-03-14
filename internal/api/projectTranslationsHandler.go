package api

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

//==========================================================================
// Endpoints
//==========================================================================

// TranslationCSVHandler - GET
func TranslationCSVHandler(gin *gin.Context) {
	handleCsvRequest(gin, "all_translations", translation.LoadTranslationsAsCSV)
}

// MissingTranslationCSVHandler - GET
func MissingTranslationCSVHandler(gin *gin.Context) {
	handleCsvRequest(gin, "missing_translations", translation.FindMissingTranslationsCSV)
}

//==========================================================================
// Private
//==========================================================================

// handleCsvRequest handles an incoming csv request. Sending the csv file produced by csvFunc back to the client.
// #### params
// gin - the gin context
// fileName - the file name to give the client
// csvFunc - a function that give a project produces csv output
func handleCsvRequest(gin *gin.Context, fileName string, csvFunc func(project *git.Project) (string, error)) {
	projectName := gin.Param("project")
	log.Logger.Infof("Incoming CSV download request for project %s", projectName)

	project := configuration.Get().GetProject(projectName)
	defer project.Unlock()

	if project == nil {
		log.Logger.Warn("Project %s not found during CSV download request", projectName)
		gin.Status(http.StatusBadRequest)
	} else {
		csvData, err := csvFunc(project)

		if err != nil {
			log.Logger.Error(err)
			gin.Status(http.StatusInternalServerError)
		} else {
			responseWriter := gin.Writer
			responseWriter.Header().Add("Content-Type", "text/csv")
			responseWriter.Header().Add(
				"Content-Disposition",
				fmt.Sprintf("attachment; filename=%s_%s_%s.csv", projectName, fileName, time.Now().Format("2006-01-02")))

			_, err = responseWriter.Write([]byte(csvData))
			if err != nil {
				log.Logger.Error("Failed to write CSV response")
			}

			responseWriter.CloseNotify()
		}
	}
}
