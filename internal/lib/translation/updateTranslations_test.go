package translation

import (
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"golang.org/x/exp/slices"
	"testing"
)

type updateTranslationTestInput struct {
	translations []translationFile.Translation
	mappings     []translationFile.Translation
	resultOk     func(translations []translationFile.Translation) bool
}

func TestApplyMappings(t *testing.T) {

	inputs := []updateTranslationTestInput{
		{
			translations: []translationFile.Translation{
				*translationFile.NewTranslation("create.me", "Create Me", []string{"french"}, map[string]string{}),
				*translationFile.NewTranslation("update.me", "Update Me", []string{"french"}, map[string]string{
					"french": "not right",
				}),
			},
			mappings: []translationFile.Translation{
				*translationFile.NewTranslation("", "Create Me", []string{"french"}, map[string]string{
					"french": "La Create Me",
				}),
				*translationFile.NewTranslation("", "Update Me", []string{"french"}, map[string]string{
					"french": "La Update Me",
				}),
			},

			resultOk: func(translations []translationFile.Translation) bool {
				idxCreate := slices.IndexFunc(translations, func(trans translationFile.Translation) bool { return trans.SourceValue == "Create Me" })
				idxUpdate := slices.IndexFunc(translations, func(trans translationFile.Translation) bool { return trans.SourceValue == "Update Me" })

				return idxCreate != -1 &&
					idxUpdate != -1 &&
					translations[idxCreate].Translations["french"] == "La Create Me" &&
					translations[idxUpdate].Translations["french"] == "La Update Me"
			},
		},
		{
			translations: []translationFile.Translation{
				*translationFile.NewTranslation("create.me", "Create Me", []string{"french"}, map[string]string{}),
				*translationFile.NewTranslation("extra.lang", "Extra language", []string{"french", "spanish"}, map[string]string{
					"french": "La Extra",
				}),
				*translationFile.NewTranslation("extra.lang.keep", "Extra language. Keep it", []string{"french", "spanish"}, map[string]string{
					"french":  "La Extra",
					"spanish": "Don't change",
				}),
				*translationFile.NewTranslation("dont.delete.me", "Keep Me", []string{"french"}, map[string]string{
					"french": "La Keep Me",
				}),
			},
			mappings: []translationFile.Translation{
				*translationFile.NewTranslation("", "Create Me", []string{"french"}, map[string]string{
					"french": "La Create Me",
				}),
				*translationFile.NewTranslation("", "Extra language", []string{"french", "spanish"}, map[string]string{
					"spanish": "Est extra languageq",
				}),
			},

			resultOk: func(translations []translationFile.Translation) bool {
				idxCreate := slices.IndexFunc(translations, func(trans translationFile.Translation) bool { return trans.SourceValue == "Create Me" })
				idxExtraUpdate := slices.IndexFunc(translations, func(trans translationFile.Translation) bool { return trans.SourceValue == "Extra language" })
				idxKeepMe := slices.IndexFunc(translations, func(trans translationFile.Translation) bool { return trans.SourceValue == "Keep Me" })
				idxExtraLangKeepMe := slices.IndexFunc(translations, func(trans translationFile.Translation) bool { return trans.SourceValue == "Extra language. Keep it" })

				return idxCreate != -1 &&
					idxExtraUpdate != -1 &&
					idxKeepMe != -1 &&
					idxExtraLangKeepMe != -1 &&
					translations[idxCreate].Translations["french"] == "La Create Me" &&
					translations[idxExtraUpdate].Translations["french"] == "La Extra" &&
					translations[idxExtraUpdate].Translations["spanish"] == "Est extra language" &&
					translations[idxKeepMe].Translations["french"] == "La Keep Me" &&
					translations[idxExtraLangKeepMe].Translations["french"] == "La Extra" &&
					translations[idxExtraLangKeepMe].Translations["spanish"] == "Don't change"
			},
		},
	}

	for i, input := range inputs {
		applyMappings(input.translations, input.mappings)

		if !input.resultOk(input.translations) {
			t.Fatalf("Error updating translations for test input %d", i)
		}
	}

}
