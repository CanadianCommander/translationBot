package translationFile

import (
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"os"
	"path"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

type PropertiesFileWriter struct {
}

func (propWriter *PropertiesFileWriter) Write(
	filePath string,
	lang string,
	sourceLanguage string,
	translations map[string]*Translation) error {

	translationsList := propWriter.sortTranslations(sourceLanguage == lang, translations)

	builder := strings.Builder{}

	for _, translation := range translationsList {
		langString, exists := translation.GetString(lang)
		if exists {
			builder.WriteString(fmt.Sprintf("%s=%s\n", translation.Key, langString))
		}
	}

	err := os.WriteFile(filePath, []byte(builder.String()), os.FileMode(0777))
	return err
}

func (propWriter *PropertiesFileWriter) CanWrite(filePath string) bool {
	return path.Ext(filePath) == ".properties"
}

//==========================================================================
// Private Methods
//==========================================================================

// sortTranslations in to the correct order they should be written in to disk
func (propWriter *PropertiesFileWriter) sortTranslations(
	isSourceLang bool,
	translations map[string]*Translation) []*Translation {

	translationsList := slices.Clone(maps.Values(translations))

	if isSourceLang {
		// sort translations
		slices.SortFunc(translationsList, func(first *Translation, second *Translation) bool {
			return first.SourceLangOrder < second.SourceLangOrder
		})
	} else {
		// sort translations
		slices.SortFunc(translationsList, func(first *Translation, second *Translation) bool {
			return strings.Compare(first.Key, second.Key) < 0
		})
	}

	return translationsList
}
