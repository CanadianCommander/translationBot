package translationFile

import "github.com/CanadianCommander/translationBot/internal/lib/git"

//==========================================================================
// Public
//==========================================================================

type Writer interface {

	// Write updated translations for lang to the specified file
	// #### params
	// pack - the language pack under which these translations are being written
	// filePath - the file to write
	// lang - the language to write to that file
	// sourceLanguage - the source language for translations in the given set
	// translations - translation data to read from
	Write(pack *git.LanguagePack, filePath string, lang string, sourceLanguage string, translations map[string]*Translation) error

	// CanWrite checks if this writer can write to the specified file
	CanWrite(filePath string) bool
}
