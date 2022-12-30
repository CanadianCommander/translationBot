package translation

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"github.com/CanadianCommander/translationBot/internal/lib/translationMapping"
	"golang.org/x/exp/maps"
)

//==========================================================================
// Public
//==========================================================================

// UpdateTranslationsFromSlackFile apply the translation updates to the specified project as detailed by the slack file.
// #### params
// slackFile - the slack file containing the update information
// project - the project to update
// #### return
// a new branch name where the updates have been applied (branched from the default branch)
func UpdateTranslationsFromSlackFile(slackFileId string, project *git.Project) (string, error) {
	loader, err := translationMapping.GetMapperForSlackFile(slackFileId)
	if err != nil {
		return "", errors.New("unexpected error when searching for file loader")
	}

	slackFileReader, err := slackutil.DownloadSlackFileById(slackFileId)
	if err != nil {
		return "", err
	}
	defer slackFileReader.Close()

	mappings, err := loader.Load(slackFileReader)
	if err != nil {
		return "", err
	}

	currentTranslations, err := LoadTranslations(project)
	if err != nil {
		return "", err
	}

	applyMappings(maps.Values(currentTranslations), mappings)

	log.Logger.Infof("%+v", currentTranslations)

	return "yo moma", nil
}

//==========================================================================
// Private
//==========================================================================

// applyMappings applies translation mappings to the given translation set.
// #### params
// translations - the translation set on which the mappings will be applied
// mappings - the mappings to apply to the translation set.
func applyMappings(translations []translationFile.Translation, mappings []translationFile.Translation) {
	mappingKeyMap := buildMappingKeyMap(mappings)
	mappingValueMap := buildMappingValueMap(mappings)

	for _, translation := range translations {
		mapping, exists := mappingKeyMap[translation.Key]
		if exists {
			translation.SourceValue = mapping.SourceValue
			mergeTranslationsWithMappings(&translation, mapping)
		} else {
			mapping, exists = mappingValueMap[translation.SourceValue]
			if exists {
				mergeTranslationsWithMappings(&translation, mapping)
			}
		}
	}

}

// buildMappingsKeyMap produces a map that indexes translation mappings by translation key
func buildMappingKeyMap(mappings []translationFile.Translation) map[string]*translationFile.Translation {
	mappingKeyMap := make(map[string]*translationFile.Translation)

	for _, mapping := range mappings {
		if mapping.Key != "" {
			mappingKeyMap[mapping.Key] = &mapping
		}
	}

	return mappingKeyMap
}

// buildMappingValueMap produces a map that indexes translation mappings by translation value
func buildMappingValueMap(mappings []translationFile.Translation) map[string]*translationFile.Translation {
	mappingValueMap := make(map[string]*translationFile.Translation)

	for _, mapping := range mappings {
		if mapping.SourceValue != "" {
			mappingValueMap[mapping.SourceValue] = &mapping
		}
	}

	return mappingValueMap
}

// mergeTranslationsWithMappings merges the translation map with the updates in the mapping.
// #### params
// translation - the translation whose translation mappings are to be updated
// mapping - the mapping (updates) to apply
func mergeTranslationsWithMappings(translation *translationFile.Translation, mapping *translationFile.Translation) {
	for lang, value := range mapping.Translations {
		translation.Translations[lang] = value
	}
}
