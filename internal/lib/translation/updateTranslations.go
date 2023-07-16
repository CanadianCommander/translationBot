package translation

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
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
	log.Logger.Infof("Updating translations for project %s using slack file %s", project.Name, slackFileId)

	loader, err := translationMapping.GetMapperForSlackFile(slackFileId)
	if err != nil {
		return "", err
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

	log.Logger.Infof("Applying translation mappings to project %s", project.Name)
	applyMappings(maps.Values(currentTranslations), mappings)

	log.Logger.Infof("Writing translation update to disk")
	newBranch, err := updateTranslationFiles(project, currentTranslations)
	if err != nil {
		return "", err
	}

	log.Logger.Infof("Translation update complete. Changes in branch %s", newBranch)
	return newBranch, nil
}

//==========================================================================
// Private
//==========================================================================

// updateTranslationFiles updates the on disk translation files for the project
// #### params
// project - the project to update
// translations - the translations to apply to the project
// #### return
// branch under which the applied changes can be found.
func updateTranslationFiles(project *git.Project, translations map[string]*translationFile.Translation) (string, error) {
	config := configuration.Get()

	newBranch := git.GenerateNewBranchName()
	err := git.SwitchBranch(project, newBranch, true)
	defer func(project *git.Project, branch string) {
		err := git.SwitchBranch(project, branch, true)
		if err != nil {
			log.Logger.Errorf("Failed to switch back to default branch for project %s."+
				" Project maybe in bad state!", project.Name)
		}
	}(project, project.Branch)
	if err != nil {
		log.Logger.Errorf("Failed to switch branch in project %s", project.Name)
		return "", err
	}

	// update translation files
	for _, pack := range project.Packs {
		for lang, transFile := range pack.TranslationFiles {
			// if source file update disabled skip file
			if !project.UpdateSourceFile && lang == project.SourceLanguage {
				continue
			}

			writer := translationFile.GetWriterForFile(transFile)

			if writer != nil {
				err = writer.Write(project.ProjectRelativePathToAbsolute(transFile), lang, project.SourceLanguage, translations)
				if err != nil {
					return "", err
				}

			} else {
				log.Logger.Errorf("Cannot  write translations to file %s. No writer matches", transFile)
			}
		}
	}

	if !config.TestMode {
		err = git.CommitAndPushChanges(project)
		if err != nil {
			log.Logger.Errorf("Failed to commit and push changes for %s", project.Name)
			return "", err
		}
	}

	return newBranch, nil
}

// applyMappings applies translation mappings to the given translation set.
// #### params
// translations - the translation set on which the mappings will be applied
// mappings - the mappings to apply to the translation set.
func applyMappings(translations []*translationFile.Translation, mappings []translationFile.Translation) {
	mappingKeyMap := buildMappingKeyMap(mappings)
	mappingValueMap := buildMappingValueMap(mappings)

	for _, translation := range translations {
		mapping, exists := mappingKeyMap[translation.Key]
		if exists {
			translation.SourceValue = mapping.SourceValue
			mergeTranslationsWithMappings(translation, mapping)
		} else {
			mapping, exists = mappingValueMap[translation.SourceValue]
			if exists {
				mergeTranslationsWithMappings(translation, mapping)
			}
		}
	}
}

// buildMappingsKeyMap produces a map that indexes translation mappings by translation key
func buildMappingKeyMap(mappings []translationFile.Translation) map[string]*translationFile.Translation {
	mappingKeyMap := make(map[string]*translationFile.Translation)

	for _, mapping := range mappings {
		if mapping.Key != "" {
			mappingCpy := mapping
			mappingKeyMap[mapping.Key] = &mappingCpy
		}
	}

	return mappingKeyMap
}

// buildMappingValueMap produces a map that indexes translation mappings by translation value
func buildMappingValueMap(mappings []translationFile.Translation) map[string]*translationFile.Translation {
	mappingValueMap := make(map[string]*translationFile.Translation)

	for idx, mapping := range mappings {
		if mapping.SourceValue != "" {
			mappingValueMap[mapping.SourceValue] = &mappings[idx]
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
		if value != "" {
			translation.Translations[lang] = value
		}
	}
}
