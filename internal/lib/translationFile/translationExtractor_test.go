package translationFile

import "testing"

type extractorTestInput struct {
	rawStructInput       map[string]interface{}
	rawStructInputFrench map[string]interface{}
	resultOk             func(tm map[string]Translation) bool
}

func TestExtractTranslations(t *testing.T) {

	inputs := []extractorTestInput{
		{
			rawStructInput: map[string]interface{}{
				"card":   "This is a card or something",
				"police": "Frank Drebin",
			},
			resultOk: func(tm map[string]Translation) bool {
				return tm["card"].SourceValue == "This is a card or something" &&
					tm["police"].SourceValue == "Frank Drebin"
			},
		},
		{
			rawStructInput: map[string]interface{}{
				"card": "This is a card or something",
				"nested": map[string]interface{}{
					"bob": "Ross",
				},
			},
			resultOk: func(tm map[string]Translation) bool {
				return tm["card"].SourceValue == "This is a card or something" &&
					tm["nested.bob"].SourceValue == "Ross"
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
			resultOk: func(tm map[string]Translation) bool {
				return tm["card"].SourceValue == "This is a card or something" &&
					tm["card"].Translations["french"] == "La This is a card or something" &&
					tm["nested.bob"].SourceValue == "Ross" &&
					tm["nested.bob"].Translations["french"] == "La Ross" &&
					tm["nested.deeper.go"].SourceValue == "deeper" &&
					tm["nested.deeper.go"].Translations["french"] == "La deeper"
			},
		},
	}

	for i, input := range inputs {
		translations := make(map[string]Translation)

		if input.rawStructInputFrench != nil {
			// reverse loading has been an issue. Load French first to test for this edge case.
			extractTranslations(
				"english",
				[]string{"french"},
				"french",
				"",
				input.rawStructInputFrench,
				translations)
		}

		extractTranslations(
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
