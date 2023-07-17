package translationFile

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"gopkg.in/yaml.v2"
	"strings"
)

//==========================================================================
// Private
//==========================================================================

// extractTranslations pulls out a list of translations from a json structure
// #### params
// pack - the language pack we are extracting translations from
// sourceLanguage - the source language from which translations are derived
// translationLanguages - the languages that can be derived from the sourceLanguage
// lang - the language of the translations being loaded
// keypath - used for recursion just set to ""
// translationData - the translation data to extract translations from
// translations - the translation map to update with the new translations
func extractTranslations(
	pack *git.LanguagePack,
	sourceLanguage string,
	translationLanguages []string,
	lang string,
	keypath string,
	translationData map[string]interface{},
	translations map[string]*Translation) {

	if keypath == "" {
		keypath = pack.Name
	}

	for key, val := range translationData {
		switch val.(type) {

		case map[string]interface{}:
			extractTranslations(
				pack,
				sourceLanguage,
				translationLanguages,
				lang,
				formatKeypath(keypath, key),
				val.(map[string]interface{}),
				translations)
		case string:
			trans, exists := translations[formatKeypath(keypath, key)]
			if exists && (trans.Pack == nil || trans.Pack.Name == pack.Name) {
				if lang == sourceLanguage {
					trans.SourceValue = val.(string)
				} else {
					trans.Translations[lang] = val.(string)
				}

				translations[formatKeypath(keypath, key)] = trans
			} else {
				if lang == sourceLanguage {
					translations[formatKeypath(keypath, key)] =
						NewTranslation(
							pack,
							stripPackNameFromKeypath(pack, formatKeypath(keypath, key)),
							strings.Trim(val.(string), " "),
							sourceLanguage,
							translationLanguages,
							make(map[string]string))
				} else {
					translationMap := map[string]string{lang: strings.Trim(val.(string), " ")}

					translations[formatKeypath(keypath, key)] =
						NewTranslation(
							pack,
							stripPackNameFromKeypath(pack, formatKeypath(keypath, key)),
							"",
							sourceLanguage,
							translationLanguages,
							translationMap)
				}
			}
		default:
			log.Logger.Warn("JSON data found that is not a string or a map! did you put a number / array in the translationData? Ignoring")
		}
	}
}

// extractTranslationsMapSlice just like extractTranslations but extracts from a MapSlice instead of a map[string]interface{}
// The big difference is that MapSlices are ordered. So when loading from them order will be stored in the produced translation objects.
// #### params
// pack - the language pack we are extracting translations from
// sourceLanguage - the source language from which translations are derived
// translationLanguages - the languages that can be derived from the sourceLanguage
// lang - the language of the translations being loaded
// keypath - used for recursion just set to ""
// translationData - the translation data to extract translations from
// translations - the translation map to update with the new translations
func extractTranslationsMapSlice(
	pack *git.LanguagePack,
	sourceLanguage string,
	translationLanguages []string,
	lang string,
	keypath string,
	translationData yaml.MapSlice,
	translations map[string]*Translation) {

	var order uint = 0
	extractTranslationsMapSliceRecursive(pack, sourceLanguage, translationLanguages, lang, keypath, translationData, translations, &order)
}

// extractTranslationsMapSliceRecursive just like extractTranslations but extracts from a MapSlice instead of a map[string]interface{}
// The big difference is that MapSlices are ordered. So when loading from them order will be stored in the produced translation objects.
// #### params
// pack - the language pack we are extracting translations from
// sourceLanguage - the source language from which translations are derived
// translationLanguages - the languages that can be derived from the sourceLanguage
// lang - the language of the translations being loaded
// keypath - used for recursion just set to ""
// translationData - the translation data to extract translations from
// translations - the translation map to update with the new translations
func extractTranslationsMapSliceRecursive(
	pack *git.LanguagePack,
	sourceLanguage string,
	translationLanguages []string,
	lang string,
	keypath string,
	translationData yaml.MapSlice,
	translations map[string]*Translation,
	order *uint) {

	if keypath == "" {
		keypath = pack.Name
	}

	for _, item := range translationData {
		switch item.Value.(type) {

		case yaml.MapSlice:
			extractTranslationsMapSliceRecursive(
				pack,
				sourceLanguage,
				translationLanguages,
				lang,
				formatKeypath(keypath, item.Key.(string)),
				item.Value.(yaml.MapSlice),
				translations,
				order)

		case string:
			trans, exists := translations[formatKeypath(keypath, item.Key.(string))]
			if exists && (trans.Pack == nil || trans.Pack.Name == pack.Name) {
				if lang == sourceLanguage {
					trans.SourceValue = item.Value.(string)
					trans.SourceLangOrder = *order
					*order++
				} else {
					trans.Translations[lang] = item.Value.(string)
				}

				translations[formatKeypath(keypath, item.Key.(string))] = trans
			} else {
				if lang == sourceLanguage {
					translations[formatKeypath(keypath, item.Key.(string))] =
						NewTranslationOrdered(
							pack,
							stripPackNameFromKeypath(pack, formatKeypath(keypath, item.Key.(string))),
							strings.Trim(item.Value.(string), " "),
							sourceLanguage,
							translationLanguages,
							make(map[string]string),
							*order)
					*order++
				} else {
					translationMap := map[string]string{lang: strings.Trim(item.Value.(string), " ")}

					translations[formatKeypath(keypath, item.Key.(string))] =
						NewTranslation(
							pack,
							stripPackNameFromKeypath(pack, formatKeypath(keypath, item.Key.(string))),
							"",
							sourceLanguage,
							translationLanguages,
							translationMap)
				}
			}
		default:
			log.Logger.Warn("JSON data found that is not a string or a map! did you put a number / array in the translationData? Ignoring")
		}
	}
}
