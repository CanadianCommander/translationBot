package translationFile

import (
	"errors"
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"regexp"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

const (
	StringNodePlaceholder = "TBStrNode"
)

type PropertiesFileLoader struct {
}

// Load the specified Properties File
// #### params
// sourceLanguage - the source language from which translations are derived
// translationLanguages - the languages that can be derived from the sourceLanguage
// language - language of the translation data
// pack - the language pack to load the translation data from
// translations - [optional] if provided translations will be added to the specified slice
// #### return
// a map of keypath: translation pairs. where keypath is the JSON path to the translation. ex. page.home.title
func (pLoader *PropertiesFileLoader) Load(
	sourceLanguage string,
	translationLanguages []string,
	language string,
	pack *git.LanguagePack,
	translations map[string]*Translation) (map[string]*Translation, error) {

	fileData, err := os.ReadFile(pack.Project.ProjectRelativePathToAbsolute(pack.TranslationFiles[language]))
	if err != nil {
		return nil, err
	}

	parsedProperties, err := pLoader.parsePropertiesFile(fileData)
	if err != nil {
		return nil, err
	}

	extractTranslationsMapSlice(pack, sourceLanguage, translationLanguages, language, "", parsedProperties, translations)
	return translations, nil
}

// CanLoad Properties files!
func (pLoader *PropertiesFileLoader) CanLoad(filePath string) bool {
	return path.Ext(filePath) == ".properties"
}

// parsePropertiesFile parses the supplied raw property file content and produces a generic
// map structure. Just like JSON or YAML.
func (pLoader *PropertiesFileLoader) parsePropertiesFile(fileData []byte) (yaml.MapSlice, error) {
	output := make([]yaml.MapItem, 0)
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
		value := strings.Trim(string(matches[2]), " ")
		output = pLoader.parsePropertiesRecursive(output, pathSegments, value)
	}

	return output, nil
}

func (pLoader *PropertiesFileLoader) parsePropertiesRecursive(
	root yaml.MapSlice,
	segments []string,
	value string) yaml.MapSlice {

	currSegment := strings.Trim(segments[0], " ")

	if len(segments) == 1 {
		return append(root, yaml.MapItem{Key: currSegment, Value: value})
	} else {

		existingMapIndex := slices.IndexFunc(root, func(item yaml.MapItem) bool {
			return item.Key == currSegment
		})

		if existingMapIndex == -1 {
			return append(root, yaml.MapItem{Key: currSegment, Value: pLoader.parsePropertiesRecursive(make(yaml.MapSlice, 0), segments[1:], value)})
		} else {

			switch root[existingMapIndex].Value.(type) {
			case string:
				root[existingMapIndex].Value = yaml.MapSlice([]yaml.MapItem{
					yaml.MapItem{Key: StringNodePlaceholder, Value: root[existingMapIndex].Value.(string)},
				})
			}

			root[existingMapIndex].Value =
				pLoader.parsePropertiesRecursive(root[existingMapIndex].Value.(yaml.MapSlice), segments[1:], value)

			return root
		}
	}
}
