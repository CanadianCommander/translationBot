package translationFile

import (
	"golang.org/x/exp/maps"
	"gopkg.in/yaml.v2"
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
func translationsToMapSlice(lang string, translationMap map[string]*Translation, isSourceLang bool) yaml.MapSlice {

	output := yaml.MapSlice{}
	translations := maps.Values(translationMap)

	if isSourceLang {
		// sort by order of appearance in original file
		sort.Slice(translations, func(i int, j int) bool {
			return translations[i].SourceLangOrder < translations[j].SourceLangOrder
		})

		// filter out translation items missing source values.
		hasSourceTranslations := make([]*Translation, 0, len(translations))
		for _, trans := range translations {
			if trans.SourceValue != "" {
				hasSourceTranslations = append(hasSourceTranslations, trans)
			}
		}
		translations = hasSourceTranslations

	} else {
		// sort alphabetically by keypath
		sort.Slice(translations, func(i int, j int) bool {
			return strings.Compare(translations[i].Key, translations[j].Key) < 0
		})
	}

	// TODO

	return output
}

// findByKeyMapSlice finds a MapItem by key in the given map slice, at the current level.
// #### params
// key - key to search for
// mapSlice - map slice to search
// #### return
// the map item if found & true / false indicating if the item was found or not
func findByKeyMapSlice(key string, mapSlice yaml.MapSlice) (*yaml.MapItem, bool) {
	for _, item := range mapSlice {
		if item.Key == key {
			return &item, true
		}
	}
	return nil, false
}
