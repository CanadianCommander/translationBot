package translationMapping

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
	"io"
)

//==========================================================================
// Public
//==========================================================================

// GetMapperForSlackFile - the appropriate translation mapper for the given slack file.
// #### params
// fileId - the id of the slack file
// #### return
// a mapping loader or none if no loaders match the input file.
func GetMapperForSlackFile(fileId string) (MappingLoader, error) {

	fileInfo, _, _, err := slackutil.Api.GetFileInfo(fileId, 0, 0)
	if err != nil {
		log.Logger.Errorf("Error getting file info for file [%s]. Error: %s", fileId, err)
		return nil, err
	}

	log.Logger.Infof("Searching for loader to handle %s file", fileInfo.Mimetype)
	for _, loader := range mappers {
		if loader.IsMimeSupported(fileInfo.Mimetype) {
			match, err := checkFileMagic(loader, fileInfo)
			if err != nil {
				log.Logger.Errorf("Unexpected error checking file magic for file %s", fileId)
				return nil, err
			}

			if match {
				return loader, nil
			}
		}
	}

	return nil, nil
}

//==========================================================================
// Private
//==========================================================================

// available mappers
var mappers = []MappingLoader{
	&CsvTaraMappingLoader{},
}

// readStartOfFile to allow for file magic checking. Reads up to 128KB
// #### params
// fileReader - the file reader to read the file data from
func readStartOfFile(fileReader io.ReadCloser) ([]byte, error) {
	fileDataStart := make([]byte, 1024*128)
	bytesRead, err := fileReader.Read(fileDataStart)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}
	if bytesRead == 0 {
		return nil, errors.New("got zero bytes when trying to read slack file")
	}

	return fileDataStart, nil
}

// checkFileMagic to see if it matches the loader
// #### params
// loader - the loader to check for match
// fileInfo - the file to check
func checkFileMagic(loader MappingLoader, fileInfo *slack.File) (bool, error) {
	fileReader, err := slackutil.DownloadSlackFile(fileInfo)
	if err != nil {
		log.Logger.Errorf("Error downloading file [%s]. Error: %s", fileInfo.ID, err)
		return false, err
	}
	defer fileReader.Close()

	fileDataStart, err := readStartOfFile(fileReader)
	if err != nil {
		log.Logger.Errorf("Error reading file [%s] from slack. Error: %s", fileInfo.ID, err)
		return false, err
	}

	return loader.IsFileSupported(fileDataStart), nil
}
