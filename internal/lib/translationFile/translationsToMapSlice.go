package translationFile

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
	"sort"
	"strings"
)

//==========================================================================
// Private
//==========================================================================

// translationsToMap creates a raw map (JSON like) structure out of the given translations for the given language
// #### params
// pack - the language pack that translations should be mapped for
// lang - the language to extract
// translationMap - the translation data to convert
// isSourceLang - if true lang will be treated as the source lang. Different extraction logic will apply.
// #### return
// a map containing a JSON like representation of the data. Each value will be either string or map[string]interface{}
func translationsToMapSlice(
	pack *git.LanguagePack,
	lang string,
	translationMap map[string]*Translation,
	isSourceLang bool) yaml.MapSlice {

	translations := maps.Values(translationMap)

	// filter translations not part of the pack
	filteredTranslations := make([]*Translation, 0, len(translations))
	for _, trans := range translations {
		if trans.Pack.Name == pack.Name {
			filteredTranslations = append(filteredTranslations, trans)
		}
	}
	translations = filteredTranslations

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

	return translationListToMapSlice(lang, isSourceLang, translations)
}

// translationListToMapSlice convert a list of translations in to a MapSlice list
// #### params
// lang - language to creat the MapSlice from
// isSourceLang - is this the source language?
// translations - the list of translations to convert
func translationListToMapSlice(lang string, isSourceLang bool, translations []*Translation) yaml.MapSlice {

	root := yaml.MapItem{
		Key:   nil,
		Value: yaml.MapSlice{},
	}

	for _, translation := range translations {
		pathParts := translation.PathParts()

		value, exists := extractValue(lang, translation, isSourceLang)
		if exists {
			mapSliceDeepCreate(&root, value, pathParts, []string{})
		}

	}

	return root.Value.(yaml.MapSlice)
}

// mapSliceDeepCreate helper method for translationsListToMapSlice. Handles the recursion down the MapSlice tree
// #### params
// item - current node in the MapSlice tree
// value - value to create
// pathParts - the path at which the value should be created
// basePath - current path location in the tree
func mapSliceDeepCreate(item *yaml.MapItem, value string, pathParts []string, basePath []string) {
	relativePathParts := toRelativePathParts(basePath, pathParts)

	if len(relativePathParts) == 1 {
		item.Value = append(item.Value.(yaml.MapSlice), yaml.MapItem{Key: relativePathParts[0], Value: value})
	} else {
		existingIdx := slices.IndexFunc(
			item.Value.(yaml.MapSlice),
			func(item yaml.MapItem) bool { return item.Key == relativePathParts[0] })

		var deepMap *yaml.MapItem = nil

		if existingIdx != -1 {
			deepMap = &item.Value.(yaml.MapSlice)[existingIdx]
		} else {
			deepMap = &yaml.MapItem{
				Key:   relativePathParts[0],
				Value: yaml.MapSlice{},
			}
		}

		mapSliceDeepCreate(deepMap, value, pathParts, pathParts[:len(pathParts)-len(relativePathParts)+1])

		if existingIdx == -1 {
			item.Value = append(item.Value.(yaml.MapSlice), *deepMap)
		}
	}
}

// toRelativePathParts makes the path parts relative to the supplied base path
// #### params
// basePath - absolute base path to be relative to
// pathParts - absolute path to make relative
func toRelativePathParts(basePath []string, pathParts []string) []string {
	outParts := make([]string, 0, len(pathParts))

	for idx, part := range pathParts {
		if idx < len(basePath) && basePath[idx] != part {
			break
		} else if idx >= len(basePath) {
			outParts = append(outParts, part)
		}
	}

	return outParts
}

// findByKeyMapSlice finds a MapItem by key in the given map slice, at the current level.
// #### params
// key - key to search for
// mapSlice - map slice to search
// #### return
// the map item if found & true / false indicating if the item was found or not
func mapSliceFindByKey(key string, mapSlice yaml.MapSlice) (*yaml.MapItem, bool) {
	for _, item := range mapSlice {
		if item.Key == key {
			return &item, true
		}
	}
	return nil, false
}
