package translationFile

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"strings"
)

//==========================================================================
// Private
//==========================================================================

// extractTranslations pulls out a list of translations from a json structure
// #### params
// sourceLanguage - the source language from which translations are derived
// translationLanguages - the languages that can be derived from the sourceLanguage
// lang - the language of the translations being loaded
// keypath - used for recursion just set to ""
// translationData - the translation data to extract translations from
// translations - the translation map to update with the new translations
func extractTranslations(
	sourceLanguage string,
	translationLanguages []string,
	lang string,
	keypath string,
	translationData map[string]interface{},
	translations map[string]Translation) {

	for key, val := range translationData {
		switch val.(type) {

		case map[string]interface{}:
			extractTranslations(
				sourceLanguage,
				translationLanguages,
				lang,
				formatKeypath(keypath, key),
				val.(map[string]interface{}),
				translations)

		case string:
			trans, exists := translations[formatKeypath(keypath, key)]
			if exists {
				trans.Translations[lang] = val.(string)
			} else {
				if lang == sourceLanguage {
					translations[formatKeypath(keypath, key)] =
						*NewTranslation(
							formatKeypath(keypath, key),
							strings.Trim(val.(string), " "),
							translationLanguages,
							make(map[string]string))
				} else {
					translationMap := map[string]string{lang: strings.Trim(val.(string), " ")}

					translations[formatKeypath(keypath, key)] =
						*NewTranslation(
							formatKeypath(keypath, key),
							"",
							translationLanguages,
							translationMap)
				}
			}
		default:
			log.Logger.Warn("JSON data found that is not a string or a map! did you put a number / array in the translationData? Ignoring")
		}
	}
}
