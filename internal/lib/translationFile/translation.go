package translationFile

import (
	"errors"
	"fmt"
	"golang.org/x/exp/maps"
)

//==========================================================================
// Public
//==========================================================================

type Translation struct {
	Key         string
	SourceValue string

	// language -> translation of source value
	Translations map[string]string

	// complete list of languages supported by this translation. Translations may not be provided for all languages.
	Languages []string
}

// NewTranslation creates a new translation object
// #### params
// key - translation file key
// sourceValue - the source value for this translation, usually the English version
// supportedLangs - all languages supported (master lanaguage list)
// translations... - one or translation mappings. lang -> translation
func NewTranslation(key string, sourceValue string, supportedLangs []string, translations map[string]string) *Translation {

	return &Translation{
		Key:          key,
		SourceValue:  sourceValue,
		Translations: translations,
		Languages:    supportedLangs,
	}
}

// Merge merges this translation with another translation
// #### params
// other - the translation to merge in to this one.
func (t *Translation) Merge(other Translation) error {
	if t.Key != other.Key {
		return errors.New(fmt.Sprintf("Attempt to merge translations with different keys! %s %s", t.Key, other.Key))
	}

	if len(other.SourceValue) != 0 {
		t.SourceValue = other.SourceValue
	}

	for lang, translation := range other.Translations {
		t.Translations[lang] = translation
	}

	allLanguages := append(t.Languages, other.Languages...)
	uniqueLanguages := map[string]string{}
	for _, lang := range allLanguages {
		uniqueLanguages[lang] = lang
	}
	t.Languages = maps.Values(uniqueLanguages)

	return nil
}

//==========================================================================
// Getters
//==========================================================================

// HumanKey formats the key of this translation for human consumption (makes it look nice)
func (t *Translation) HumanKey() string {
	return unEscapeKeypath(t.Key)
}

// PathParts gets the path parts for this translation's key
func (t *Translation) PathParts() []string {
	return splitKeypath(t.Key)
}

// MissingLanguages produces a list of languages missing translations
func (t *Translation) MissingLanguages() []string {
	var missing []string

	// map reduce filter. Why go no like them :(
	for _, lang := range t.Languages {
		if _, present := t.Translations[lang]; !present {
			missing = append(missing, lang)
		}
	}

	return missing
}

func (t *Translation) HasMissingTranslations() bool {
	return len(t.MissingLanguages()) != 0
}
