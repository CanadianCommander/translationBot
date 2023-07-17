package translationFile

import "github.com/CanadianCommander/translationBot/internal/lib/git"

//==========================================================================
// Public
//==========================================================================

type Loader interface {
	// Load the specified translation data
	// #### params
	// sourceLanguage - the source language from which translations are derived
	// translationLanguages - the languages that can be derived from the sourceLanguage
	// language - language of the translation data
	// file - path to file containing translation data
	// translations - [optional] if provided translations will be added to the specified slice
	// #### return
	// a map of keypath: translation pairs. where keypath is the JSON path to the translation. ex. page.home.title
	Load(sourceLanguage string,
		translationLanguages []string,
		language string,
		pack *git.LanguagePack,
		translations map[string]*Translation,
	) (map[string]*Translation, error)

	// CanLoad checks if this loader can load the specified file
	CanLoad(filePath string) bool
}
