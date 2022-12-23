package translationMapping

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"mime"
	"regexp"
)

//==========================================================================
// Public
//==========================================================================

const CsvTaraFileMagic = "^ENGLISH,\\s+FRENCH"

type CsvTaraMappingLoader struct {
}

func (c *CsvTaraMappingLoader) Load(fileData []byte) []map[string]string {
	//TODO implement me
	panic("implement me")
}

func (c *CsvTaraMappingLoader) IsMimeSupported(mimetype string) bool {
	parsedMime, _, err := mime.ParseMediaType(mimetype)
	if err != nil {
		log.Logger.Error("Mime type parsing error ", err)
		return false
	}

	csvMime, _, _ := mime.ParseMediaType(mime.TypeByExtension(".csv"))
	return parsedMime == csvMime
}

func (c *CsvTaraMappingLoader) IsFileSupported(fileDataStart []byte) bool {
	matched, err := regexp.Match(CsvTaraFileMagic, fileDataStart)
	return matched && err == nil
}

func (c *CsvTaraMappingLoader) GetLoaderType() string {
	return "CSV_Tara_1.0"
}
