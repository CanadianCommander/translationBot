package translation

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"golang.org/x/exp/slices"
	"testing"
)

type updateTranslationTestInput struct {
	translations []*translationFile.Translation
	mappings     []translationFile.Translation
	resultOk     func(translations []*translationFile.Translation) bool
}

func TestApplyMappings(t *testing.T) {

	dummyPackage := git.LanguagePack{
		Name: "dummy",
		TranslationFiles: map[string]string{
			"english": "english.json",
			"french":  "french.json",
		},
		Project: &git.Project{},
	}

	inputs := []updateTranslationTestInput{
		{
			translations: []*translationFile.Translation{
				translationFile.NewTranslation(&dummyPackage, "create.me", "Create Me", "english", []string{"french"}, map[string]string{}),
				translationFile.NewTranslation(&dummyPackage, "update.me", "Update Me", "english", []string{"french"}, map[string]string{
					"french": "not right",
				}),
			},
			mappings: []translationFile.Translation{
				*translationFile.NewTranslation(&dummyPackage, "", "Create Me", "english", []string{"french"}, map[string]string{
					"french": "La Create Me",
				}),
				*translationFile.NewTranslation(&dummyPackage, "", "Update Me", "english", []string{"french"}, map[string]string{
					"french": "La Update Me",
				}),
			},

			resultOk: func(translations []*translationFile.Translation) bool {
				idxCreate := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.SourceValue == "Create Me" })
				idxUpdate := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.SourceValue == "Update Me" })

				return idxCreate != -1 &&
					idxUpdate != -1 &&
					translations[idxCreate].Translations["french"] == "La Create Me" &&
					translations[idxUpdate].Translations["french"] == "La Update Me"
			},
		},
		{
			translations: []*translationFile.Translation{
				translationFile.NewTranslation(&dummyPackage, "create.me", "Create Me", "english", []string{"french"}, map[string]string{}),
				translationFile.NewTranslation(&dummyPackage, "extra.lang", "Extra language", "english", []string{"french", "spanish"}, map[string]string{
					"french": "La Extra",
				}),
				translationFile.NewTranslation(&dummyPackage, "extra.lang.keep", "Extra language. Keep it", "english", []string{"french", "spanish"}, map[string]string{
					"french":  "La Extra",
					"spanish": "Don't change",
				}),
				translationFile.NewTranslation(&dummyPackage, "dont.delete.me", "Keep Me", "english", []string{"french"}, map[string]string{
					"french": "La Keep Me",
				}),
			},
			mappings: []translationFile.Translation{
				*translationFile.NewTranslation(&dummyPackage, "", "Create Me", "english", []string{"french"}, map[string]string{
					"french": "La Create Me",
				}),
				*translationFile.NewTranslation(&dummyPackage, "", "Extra language", "english", []string{"french", "spanish"}, map[string]string{
					"spanish": "Est extra language",
				}),
			},

			resultOk: func(translations []*translationFile.Translation) bool {
				idxCreate := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.SourceValue == "Create Me" })
				idxExtraUpdate := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.SourceValue == "Extra language" })
				idxKeepMe := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.SourceValue == "Keep Me" })
				idxExtraLangKeepMe := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.SourceValue == "Extra language. Keep it" })

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
		{
			translations: []*translationFile.Translation{
				translationFile.NewTranslation(&dummyPackage, "foo.zip", "Zip", "english", []string{"french"}, map[string]string{
					"spanish": "El Zip",
				}),
				translationFile.NewTranslation(&dummyPackage, "foo.bang", "Bang", "english", []string{"spanish"}, map[string]string{
					"french": "La Bang",
				}),
			},

			mappings: []translationFile.Translation{
				*translationFile.NewTranslation(&dummyPackage, "foo.zip", "Zip", "english", []string{"french"}, map[string]string{
					"french": "La Zip",
				}),
				*translationFile.NewTranslation(&dummyPackage, "foo.bang", "Bang", "english", []string{"spanish"}, map[string]string{
					"spanish": "El Bang",
				}),
				*translationFile.NewTranslation(&dummyPackage, "foo.banger", "Bang", "english", []string{"spanish"}, map[string]string{
					"spanish": "No",
				}),
			},

			resultOk: func(translations []*translationFile.Translation) bool {
				idxZip := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.Key == "foo.zip" })
				idxBang := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.Key == "foo.bang" })
				idxBanger := slices.IndexFunc(translations, func(trans *translationFile.Translation) bool { return trans.Key == "foo.banger" })

				return idxZip != -1 &&
					idxBang != -1 &&
					idxBanger == -1 &&
					translations[idxZip].Translations["french"] == "La Zip" &&
					translations[idxBang].Translations["spanish"] == "El Bang"
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
