package translation

import "github.com/CanadianCommander/translationBot/internal/lib/git"

//==========================================================================
// Public
//==========================================================================

// LoadTranslations - load translations for the given project
// #### params
// project - the project who's translations will be loaded
func LoadTranslations(project *git.Project) ([]Translation, error) {

	err := git.PullProjectRepo(project)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
