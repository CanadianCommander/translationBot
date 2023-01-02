package translationFile

import (
	"golang.org/x/exp/maps"
	"sort"
	"strings"
)

//==========================================================================
// Private
//==========================================================================

// translationsToMap creates a raw map (JSON like) structure out of the given translations for the given language
// #### params
// lang - the language to extract
// translationMap - the translation data to convert
// isSourceLang - if true lang will be treated as the source lang. Different extraction logic will apply.
// #### return
// a map containing a JSON like representation of the data. Each value will be either string or map[string]interface{}
func translationsToMap(lang string, translationMap map[string]Translation, isSourceLang bool) map[string]interface{} {

	output := make(map[string]interface{})
	translations := maps.Values(translationMap)

	// sort alphabetically by keypath
	sort.Slice(translations, func(i int, j int) bool {
		return strings.Compare(translations[i].Key, translations[j].Key) < 0
	})

	for _, translation := range translations {
		pathParts := strings.Split(translation.Key, ".")
		currentPosition := output

		for i, part := range pathParts {
			if i == len(pathParts)-1 {
				val, exists := extractValue(lang, translation, isSourceLang)
				if exists {
					currentPosition[part] = val
				}
			} else {
				lowerMap, exists := currentPosition[part]
				if exists {
					currentPosition = lowerMap.(map[string]interface{})
				} else {
					currentPosition[part] = make(map[string]interface{})
					currentPosition = currentPosition[part].(map[string]interface{})
				}
			}
		}
	}

	return output
}

// extractValue from the translation for the given lang. Returning the value and if that value exists or not.
func extractValue(lang string, translation Translation, isSourceLang bool) (string, bool) {
	if isSourceLang {
		return translation.SourceValue, true
	} else {
		val, exists := translation.Translations[lang]
		return val, exists
	}
}
