package slackutil

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/slack-go/slack"
	"io"
	"net/http"
)

// DownloadSlackFileById
// #### params
// fileId - the id of the slack file to download
// #### return
// io reader containing file data
func DownloadSlackFileById(fileId string) (io.ReadCloser, error) {
	fileInfo, _, _, err := Api.GetFileInfo(fileId, 0, 0)
	if err != nil {
		log.Logger.Errorf("Error getting file info for file [%s]. Error: %s", fileId, err)
		return nil, err
	}

	return DownloadSlackFile(fileInfo)
}

// DownloadSlackFile
// #### params
// slackFile - the slack file to download
// #### return
// io reader containing file data
func DownloadSlackFile(slackFile *slack.File) (io.ReadCloser, error) {
	log.Logger.Infof("Downloading file from slack %s", slackFile.URLPrivateDownload)

	// Custom download logic instead of using slack library. This is to allow for streaming of file content.
	// whereas the library forces buffering of the entire file in memory.
	request, _ := http.NewRequest("GET", slackFile.URLPrivateDownload, nil)
	request.Header.Add("Authorization", "Bearer "+configuration.Get().SlackClientSecret)
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Logger.Infof("Failed to download file [%s] from slack with error %s", slackFile.URLPrivateDownload, err)
		return nil, err
	}

	return response.Body, nil
}
