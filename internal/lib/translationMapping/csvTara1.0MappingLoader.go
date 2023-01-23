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

const CsvTaraFileMagic = "^\\s*ENGLISH\\s*,\\s*FRENCH"

type CsvTaraMappingLoader struct {
}

func (c *CsvTaraMappingLoader) Load(fileData io.Reader) ([]translationFile.Translation, error) {
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
		} else if lenNonBlank(mapping) < len(headers) {
			return nil, NewValidationError("one or more CSV rows have less columns then there are headers")
		}

		translationMap := make(map[string]string)
		for idx, lang := range headers[1:] {
			translationMap[lang] = strings.Trim(mapping[idx+1], " ")
		}

		mappings = append(
			mappings,
			*translationFile.NewTranslation(
				"",
				strings.Trim(mapping[0], " "),
				headers[1:],
				translationMap),
		)
	}

	return mappings, nil
}

func (c *CsvTaraMappingLoader) IsMimeSupported(mimetype string) bool {
	parsedMime, _, err := mime.ParseMediaType(mimetype)
	if err != nil {
		log.Logger.Error("Mime type parsing error ", err)
		return false
	}

	return parsedMime == "text/csv"
}

func (c *CsvTaraMappingLoader) IsFileSupported(fileDataStart []byte) bool {
	matched, err := regexp.Match(CsvTaraFileMagic, fileDataStart)
	return matched && err == nil
}

func (c *CsvTaraMappingLoader) GetLoaderType() string {
	return "CSV_Tara_1.0"
}

//==========================================================================
// Private
//==========================================================================

// lenNonBlank counts the number of non blank strings in the slice
func lenNonBlank(strings []string) int {
	len := 0
	for _, str := range strings {
		if str != "" {
			len++
		}
	}

	return len
}
