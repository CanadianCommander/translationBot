package translation

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"golang.org/x/exp/maps"
)

//==========================================================================
// Public
//==========================================================================

// FindMissingTranslations finds missing translations in the given project
// #### params
// project - the project to search for missing translations.
// #### return
// missing translation list. Where each translation is guaranteed to have a missing translation in at least one language.
func FindMissingTranslations(project *git.Project) ([]Translation, error) {
	allTranslationsMap, err := LoadTranslations(project)
	if err != nil {
		return nil, err
	}

	allTranslations := maps.Values(allTranslationsMap)
	missingTranslations := make([]Translation, 0, len(allTranslations))
	for _, trans := range allTranslations {
		if trans.HasMissingTranslations() {
			missingTranslations = append(missingTranslations, trans)
		}
	}

	return missingTranslations, nil
}
