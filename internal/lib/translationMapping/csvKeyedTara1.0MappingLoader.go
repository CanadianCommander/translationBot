package translationMapping

import (
	"encoding/csv"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"io"
	"mime"
	"regexp"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

const CsvKeyedTaraFileMagic = "^\\s*KEY\\s*,\\s*ENGLISH\\s*,\\s*FRENCH"

type CsvKeyedTaraMappingLoader struct {
}

func (c *CsvKeyedTaraMappingLoader) Load(fileData io.Reader) ([]translationFile.Translation, error) {
	csvReader := csv.NewReader(fileData)
	headers, err := csvReader.Read()
	if err != nil {
		log.Logger.Error("Error while reading csv header ", err)
		return nil, err
	} else if len(headers) < 2 {
		return nil, NewValidationError("expected at least 2 headers in CSV")
	}

	for idx, _ := range headers {
		headers[idx] = strings.Trim(strings.ToLower(headers[idx]), " ")
	}

	mappings := make([]translationFile.Translation, 0, 1024)
	for mapping, err := csvReader.Read(); mapping != nil; mapping, err = csvReader.Read() {
		if err != nil {
			return nil, err
		}

		translationMap := make(map[string]string)
		for idx, lang := range headers[2:] {
			if idx >= len(mapping) {
				translationMap[lang] = ""
			} else {
				translationMap[lang] = strings.Trim(mapping[idx+2], " ")
			}
		}

		mappings = append(
			mappings,
			*translationFile.NewTranslation(
				strings.Trim(mapping[0], " "),
				strings.Trim(mapping[1], " "),
				headers[1],
				headers[1:],
				translationMap),
		)
	}

	return mappings, nil
}

func (c *CsvKeyedTaraMappingLoader) IsMimeSupported(mimetype string) bool {
	parsedMime, _, err := mime.ParseMediaType(mimetype)
	if err != nil {
		log.Logger.Error("Mime type parsing error ", err)
		return false
	}

	return parsedMime == "text/csv"
}

func (c *CsvKeyedTaraMappingLoader) IsFileSupported(fileDataStart []byte) bool {
	matched, err := regexp.Match(CsvKeyedTaraFileMagic, fileDataStart)
	return matched && err == nil
}

func (c *CsvKeyedTaraMappingLoader) GetLoaderType() string {
	return "CSV_Keyed_Tara_1.0"
}
