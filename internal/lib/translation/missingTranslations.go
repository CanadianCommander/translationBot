package translation

import "github.com/CanadianCommander/translationBot/internal/lib/git"

//==========================================================================
// Public
//==========================================================================

// FindMissingTranslations finds missing translations in the given project
// #### params
// project - the project to search for missing translations.
// #### return
// missing translation list. Where each translation is guaranteed to have a missing translation in at least one language.
func FindMissingTranslations(project *git.Project) ([]Translation, error) {
	allTranslations, err := LoadTranslations(project)
	if err != nil {
		return nil, err
	}
	return allTranslations, nil
}
