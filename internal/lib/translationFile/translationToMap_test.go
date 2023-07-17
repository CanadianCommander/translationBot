package translationFile

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"testing"
)

type mapperTestInput struct {
	translations map[string]*Translation
	englishOk    func(rawMap map[string]interface{}) bool
	frenchOk     func(rawMap map[string]interface{}) bool
}

func TestTranslationsToMap(t *testing.T) {

	dummyPack := git.LanguagePack{
		Name: "dummy",
	}

	inputs := []mapperTestInput{
		{
			translations: map[string]*Translation{
				"foo": {
					Key:         "foo",
					Pack:        &dummyPack,
					SourceValue: "Bar",
					Translations: map[string]string{
						"french": "La Bar",
					},
					Languages: []string{"french"},
				},
				"bing.bang.boom": {
					Key:         "bing.bang.boom",
					Pack:        &dummyPack,
					SourceValue: "Bop",
					Translations: map[string]string{
						"french": "La Bop",
					},
					Languages: []string{"french"},
				},
			},
			englishOk: func(rawMap map[string]interface{}) bool {
				foo, fooExists := rawMap["foo"]
				bingBang, bingBangExists := rawMap["bing"].(map[string]interface{})["bang"].(map[string]interface{})["boom"]

				return fooExists &&
					bingBangExists &&
					foo == "Bar" &&
					bingBang == "Bop"
			},
			frenchOk: func(rawMap map[string]interface{}) bool {
				foo, fooExists := rawMap["foo"]
				bingBang, bingBangExists := rawMap["bing"].(map[string]interface{})["bang"].(map[string]interface{})["boom"]

				return fooExists &&
					bingBangExists &&
					foo == "La Bar" &&
					bingBang == "La Bop"
			},
		},
	}

	for i, input := range inputs {
		if !input.englishOk(translationsToMap(&dummyPack, "english", input.translations, true)) {
			t.Fatalf("Failed to transform translations in to map of english for input %d", i)
		}
		if !input.frenchOk(translationsToMap(&dummyPack, "french", input.translations, false)) {
			t.Fatalf("Failed to transform translations in to map of french for input %d", i)
		}
	}
}
