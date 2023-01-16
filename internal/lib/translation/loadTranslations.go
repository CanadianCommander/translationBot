package translation

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"sync"
)

//==========================================================================
// Public
//==========================================================================

// LoadTranslations - load translations for the given project
// #### params
// project - the project whose translations will be loaded
func LoadTranslations(project *git.Project) (map[string]translationFile.Translation, error) {

	err := git.InitializeProjectRepo(project)
	if err != nil {
		return nil, err
	}

	translationOutputs := make(chan map[string]translationFile.Translation, len(project.TranslationFiles))
	errorOutputs := make(chan error, len(project.TranslationFiles))
	workGroup := sync.WaitGroup{}

	for lang, file := range project.TranslationFiles {
		log.Logger.Infof("Loading translation file %s for project %s", file, project.Name)
		workGroup.Add(1)

		go translationLoadRoutine(
			project.ProjectRelativePathToAbsolute(file),
			lang,
			project.SourceLanguage,
			project.TranslationLanguages(),
			&workGroup,
			translationOutputs,
			errorOutputs,
		)
	}
	workGroup.Wait()
	close(errorOutputs)
	close(translationOutputs)

	// check for errors
	for err := range errorOutputs {
		return nil, err
	}

	// combine outputs
	var allTranslations []map[string]translationFile.Translation
	for translations := range translationOutputs {
		allTranslations = append(allTranslations, translations)
	}

	return combineTranslations(allTranslations)
}

// combineTranslations combines multiple translation maps in to one.
// #### params
// translations - translation list to combine
func combineTranslations(translations []map[string]translationFile.Translation) (map[string]translationFile.Translation, error) {
	combinedTranslations := map[string]translationFile.Translation{}

	for _, translationMap := range translations {
		for key, trans := range translationMap {
			combinedTranslation, exists := combinedTranslations[key]

			if exists {
				err := combinedTranslation.Merge(trans)
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
// file - file to load
// lang - language to load
// sourceLanguage - source language for the project
// translationLanguages - list of all translated languages for the project.
// workGroup - for this go routine
// output - channel on which the loaded translations will be returned
// errors - error channel
func translationLoadRoutine(
	file string,
	lang string,
	sourceLang string,
	translationLanguages []string,
	workGroup *sync.WaitGroup,
	output chan<- map[string]translationFile.Translation,
	errChan chan<- error) {
	defer workGroup.Done()
	translations := make(map[string]translationFile.Translation)

	loader := translationFile.GetLoaderForFile(file)
	if loader == nil {
		errChan <- errors.New("Could not find translation loader to handle file " + file)
	}

	_, err := loader.Load(
		sourceLang,
		translationLanguages,
		lang,
		file,
		translations)
	if err != nil {
		errChan <- err
	}

	output <- translations
}
