package translationFile

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"gopkg.in/yaml.v2"
	"testing"
)

type extractorTestInput struct {
	rawStructInput       map[string]interface{}
	rawStructInputFrench map[string]interface{}
	resultOk             func(tm map[string]*Translation) bool
}

func TestExtractTranslations(t *testing.T) {

	packDummy := git.LanguagePack{ // dummy pack
		Name: "dummy",
		Project: &git.Project{
			Name: "dummy",
		},
		TranslationFiles: map[string]string{
			"en": "en.json",
			"fr": "fr.json",
		},
	}

	inputs := []extractorTestInput{
		{
			rawStructInput: map[string]interface{}{
				"card":   "This is a card or something",
				"police": "Frank Drebin",
			},
			resultOk: func(tm map[string]*Translation) bool {
				return tm["dummy.card"].SourceValue == "This is a card or something" &&
					tm["dummy.police"].SourceValue == "Frank Drebin"
			},
		},
		{
			rawStructInput: map[string]interface{}{
				"card": "This is a card or something",
				"nested": map[string]interface{}{
					"bob": "Ross",
				},
			},
			resultOk: func(tm map[string]*Translation) bool {
				return tm["dummy.card"].SourceValue == "This is a card or something" &&
					tm["dummy.nested.bob"].SourceValue == "Ross"
			},
		},
		{
			rawStructInput: map[string]interface{}{
				"card": "This is a card or something",
				"nested": map[string]interface{}{
					"bob": "Ross",
					"deeper": map[string]interface{}{
						"go": "deeper",
					},
				},
			},
			rawStructInputFrench: map[string]interface{}{
				"card": "La This is a card or something",
				"nested": map[string]interface{}{
					"bob": "La Ross",
					"deeper": map[string]interface{}{
						"go": "La deeper",
					},
				},
			},
			resultOk: func(tm map[string]*Translation) bool {
				return tm["dummy.card"].SourceValue == "This is a card or something" &&
					tm["dummy.card"].Translations["french"] == "La This is a card or something" &&
					tm["dummy.nested.bob"].SourceValue == "Ross" &&
					tm["dummy.nested.bob"].Translations["french"] == "La Ross" &&
					tm["dummy.nested.deeper.go"].SourceValue == "deeper" &&
					tm["dummy.nested.deeper.go"].Translations["french"] == "La deeper"
			},
		},
	}

	for i, input := range inputs {
		translations := make(map[string]*Translation)

		if input.rawStructInputFrench != nil {
			// reverse loading has been an issue. Load French first to test for this edge case.
			extractTranslations(
				&packDummy,
				"english",
				[]string{"french"},
				"french",
				"",
				input.rawStructInputFrench,
				translations)
		}

		extractTranslations(
			&packDummy,
			"english",
			[]string{"french"},
			"english",
			"",
			input.rawStructInput,
			translations)

		if !input.resultOk(translations) {
			t.Fatalf("Failed to extract translations from input data for test set %d", i)
		}
	}
}

type extractorMapSliceTestInput struct {
	rawStructInput       yaml.MapSlice
	rawStructInputFrench yaml.MapSlice
	resultOk             func(tm map[string]*Translation) bool
}

func TestExtractTranslationsMapSlice(t *testing.T) {

	packDummy := git.LanguagePack{ // dummy pack
		Name: "dummy",
		Project: &git.Project{
			Name: "dummy",
		},
		TranslationFiles: map[string]string{
			"en": "en.json",
			"fr": "fr.json",
		},
	}

	inputs := []extractorMapSliceTestInput{
		{
			rawStructInput: yaml.MapSlice{
				{
					Key:   "foo",
					Value: "bar",
				},
			},
			rawStructInputFrench: yaml.MapSlice{
				{
					Key:   "foo",
					Value: "la bar",
				},
			},
			resultOk: func(tm map[string]*Translation) bool {
				return tm["dummy.foo"].SourceValue == "bar" &&
					tm["dummy.foo"].Translations["french"] == "la bar"
			},
		},
		{
			rawStructInput: yaml.MapSlice{
				{
					Key:   "z",
					Value: "a",
				},
				{
					Key:   "a",
					Value: "z",
				},
				{
					Key: "b",
					Value: yaml.MapSlice{
						{
							Key:   "g",
							Value: "ziziiziz",
						},
						{
							Key:   "a",
							Value: "bla",
						},
					},
				},
				{
					Key: "g",
					Value: yaml.MapSlice{
						{
							Key:   "z",
							Value: "ziziiziz",
						},
						{
							Key:   "x",
							Value: "bla",
						},
					},
				},
			},
			rawStructInputFrench: yaml.MapSlice{
				{
					Key:   "z",
					Value: "la bar",
				},
			},
			resultOk: func(tm map[string]*Translation) bool {
				return tm["dummy.z"].SourceLangOrder == 0 &&
					tm["dummy.a"].SourceLangOrder == 1 &&
					tm["dummy.b.g"].SourceLangOrder == 2 &&
					tm["dummy.b.a"].SourceLangOrder == 3 &&
					tm["dummy.g.z"].SourceLangOrder == 4 &&
					tm["dummy.g.x"].SourceLangOrder == 5
			},
		},
	}

	for i, input := range inputs {
		translations := make(map[string]*Translation)

		if input.rawStructInputFrench != nil {
			// reverse loading has been an issue. Load French first to test for this edge case.
			extractTranslationsMapSlice(
				&packDummy,
				"english",
				[]string{"french"},
				"french",
				"",
				input.rawStructInputFrench,
				translations)
		}

		extractTranslationsMapSlice(
			&packDummy,
			"english",
			[]string{"french"},
			"english",
			"",
			input.rawStructInput,
			translations)

		if !input.resultOk(translations) {
			t.Fatalf("Failed to extract translations from input data for test set %d", i)
		}
	}
}
