package translation

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"golang.org/x/exp/maps"
	"sync"
)

//==========================================================================
// Public
//==========================================================================

// LoadTranslations - load translations for the given project
// #### params
// project - the project whose translations will be loaded
func LoadTranslations(project *git.Project) (map[string]*translationFile.Translation, error) {

	err := git.InitializeProjectRepo(project)
	if err != nil {
		return nil, err
	}

	translationOutputs := make(chan map[string]*translationFile.Translation, project.TranslationFileCount())
	errorOutputs := make(chan error, project.TranslationFileCount())
	workGroup := sync.WaitGroup{}

	for _, pack := range project.Packs {
		for lang, file := range pack.TranslationFiles {
			log.Logger.Infof("Loading translation file %s for project %s", file, project.Name)
			workGroup.Add(1)

			go translationLoadRoutine(
				pack,
				lang,
				project.SourceLanguage,
				project.TranslationLanguages(),
				&workGroup,
				translationOutputs,
				errorOutputs,
			)
		}
	}
	workGroup.Wait()
	close(errorOutputs)
	close(translationOutputs)

	// check for errors
	for err := range errorOutputs {
		return nil, err
	}

	// combine outputs
	var allTranslations []map[string]*translationFile.Translation
	for translations := range translationOutputs {
		allTranslations = append(allTranslations, translations)
	}

	return combineTranslations(allTranslations)
}

// LoadTranslationsAsCSV loads all translations and outputs them as a CSV.
// #### params
// project - the project to search for missing translations.
// #### return
// CSV translation list.
func LoadTranslationsAsCSV(project *git.Project) (string, error) {
	translations, err := LoadTranslations(project)
	if err != nil {
		return "", err
	}

	return ToCSV(maps.Values(translations), project.SourceLanguage)
}

// combineTranslations combines multiple translation maps in to one.
// #### params
// translations - translation list to combine
func combineTranslations(translations []map[string]*translationFile.Translation) (map[string]*translationFile.Translation, error) {
	combinedTranslations := map[string]*translationFile.Translation{}

	for _, translationMap := range translations {
		for key, trans := range translationMap {
			combinedTranslation, exists := combinedTranslations[key]

			if exists {
				err := combinedTranslation.Merge(*trans)
				if err != nil {
					return nil, err
				}
				combinedTranslations[key] = combinedTranslation
			} else {
				combinedTranslations[key] = trans
			}
		}
	}

	return combinedTranslations, nil
}

// translationLoadRoutine loads a translation file in a go routine
// #### params
// pack - language pack to load
// lang - language to load
// sourceLanguage - source language for the project
// translationLanguages - list of all translated languages for the project.
// workGroup - for this go routine
// output - channel on which the loaded translations will be returned
// errors - error channel
func translationLoadRoutine(
	pack *git.LanguagePack,
	lang string,
	sourceLang string,
	translationLanguages []string,
	workGroup *sync.WaitGroup,
	output chan<- map[string]*translationFile.Translation,
	errChan chan<- error) {
	defer workGroup.Done()
	translations := make(map[string]*translationFile.Translation)

	file := pack.Project.ProjectRelativePathToAbsolute(pack.TranslationFiles[lang])
	loader := translationFile.GetLoaderForFile(file)
	if loader == nil {
		errChan <- errors.New("Could not find translation loader to handle file " + file)
	} else {

		_, err := loader.Load(
			sourceLang,
			translationLanguages,
			lang,
			pack,
			translations)
		if err != nil {
			errChan <- err
		}

		output <- translations
	}
}
