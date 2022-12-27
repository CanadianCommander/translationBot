package translation

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
)

//==========================================================================
// Public
//==========================================================================

// LoadTranslations - load translations for the given project
// #### params
// project - the project whose translations will be loaded
func LoadTranslations(project *git.Project) (map[string]Translation, error) {

	err := git.PullProjectRepo(project)
	if err != nil {
		return nil, err
	}

	var translations = map[string]Translation{}
	for lang, file := range project.TranslationFiles {
		log.Logger.Infof("Loading translation file %s for project %s", file, project.Name)

		loader := GetLoaderForFile(file)
		if loader == nil {
			return nil, errors.New("Could not find translation loader to handle file " + file)
		}

		_, err := loader.Load(
			project.SourceLanguage,
			project.TranslationLanguages(),
			lang,
			project.ProjectRelativePathToAbsolute(file),
			translations)
		if err != nil {
			return nil, err
		}
	}

	return translations, nil
}
