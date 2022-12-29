package translationMapping

import (
	"encoding/csv"
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"io"
	"mime"
	"regexp"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

const CsvTaraFileMagic = "^ENGLISH,\\s+FRENCH"

type CsvTaraMappingLoader struct {
}

func (c *CsvTaraMappingLoader) Load(fileData io.Reader) ([]translation.Translation, error) {
	csvReader := csv.NewReader(fileData)
	headers, err := csvReader.Read()
	if err != nil {
		log.Logger.Error("Error while reading csv header ", err)
		return nil, err
	} else if len(headers) < 2 {
		return nil, errors.New("expected at least 2 headers in CSV")
	}

	for idx, _ := range headers {
		headers[idx] = strings.Trim(strings.ToLower(headers[idx]), " ")
	}

	mappings := make([]translation.Translation, 0, 1024)
	for mapping, err := csvReader.Read(); mapping != nil; mapping, err = csvReader.Read() {
		if err != nil {
			return nil, err
		} else if len(mapping) < len(headers) {
			return nil, errors.New("one ore more CSV rows have less columns than there are headers")
		}

		translationMap := make(map[string]string)
		for idx, lang := range headers[1:] {
			translationMap[lang] = mapping[idx+1]
		}

		mappings = append(
			mappings,
			*translation.NewTranslation(
				"",
				mapping[0],
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
