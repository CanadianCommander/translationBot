package translationFile

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

type PropertiesFileLoader struct {
}

// Load the specified Properties File
// #### params
// sourceLanguage - the source language from which translations are derived
// translationLanguages - the languages that can be derived from the sourceLanguage
// language - language of the translation data
// file - path to file containing translation data
// translations - [optional] if provided translations will be added to the specified slice
// #### return
// a map of keypath: translation pairs. where keypath is the JSON path to the translation. ex. page.home.title
func (pLoader PropertiesFileLoader) Load(
	sourceLanguage string,
	translationLanguages []string,
	language string,
	file string,
	translations map[string]*Translation) (map[string]*Translation, error) {

	fileData, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	parsedProperties, err := pLoader.parsePropertiesFile(fileData)
	if err != nil {
		return nil, err
	}

	extractTranslations(sourceLanguage, translationLanguages, language, "", parsedProperties, translations)
	return translations, nil
}

// CanLoad Properties files!
func (pLoader PropertiesFileLoader) CanLoad(filePath string) bool {
	return path.Ext(filePath) == ".properties"
}

// parsePropertiesFile parses the supplied raw property file content and produces a generic
// map structure. Just like JSON or YAML.
func (pLoader PropertiesFileLoader) parsePropertiesFile(fileData []byte) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	propertyParseExpression := regexp.MustCompile("([^=]+)=(.*)")

	fileDataNoComments := regexp.MustCompile("#.*").ReplaceAllString(string(fileData[:]), "")
	lines := strings.Split(fileDataNoComments, "\n")

	for _, property := range lines {
		if len(strings.Trim(property, " ")) == 0 {
			continue
		}

		matches := propertyParseExpression.FindSubmatch([]byte(property))
		if len(matches) != 3 {
			return nil, errors.New(fmt.Sprintf("Line [%s] is malformed in properties file", property))
		}

		pathSegments := strings.Split(string(matches[1]), ".")
		currentPosition := output

		for idx, segment := range pathSegments {
			segment = strings.Trim(segment, " ")

			if idx == len(pathSegments)-1 {
				currentPosition[segment] = string(matches[2])
				break
			} else if currentPosition[segment] == nil {
				currentPosition[segment] = make(map[string]interface{})
			}
			currentPosition = currentPosition[segment].(map[string]interface{})
		}

	}

	return output, nil
}
