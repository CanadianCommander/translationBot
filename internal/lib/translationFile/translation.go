package translationFile

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

//==========================================================================
// Getters
//==========================================================================

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
