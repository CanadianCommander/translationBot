package translationFile

import (
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v2"
	"testing"
)

type TranslationsToMapSliceTestInput struct {
	translations map[string]*Translation
	englishOk    func(rawMap yaml.MapSlice) bool
	frenchOk     func(rawMap yaml.MapSlice) bool
}

func TestTranslationsToMapSlice(t *testing.T) {

	inputs := []TranslationsToMapSliceTestInput{
		{
			translations: map[string]*Translation{
				"foo": {
					Key:         "foo",
					SourceValue: "Bar",
					Translations: map[string]string{
						"french": "La Bar",
					},
					Languages: []string{"french"},
				},
				"bing.bang.boom": {
					Key:         "bing.bang.boom",
					SourceValue: "Bop",
					Translations: map[string]string{
						"french": "La Bop",
					},
					Languages: []string{"french"},
				},
			},
			englishOk: func(mapSlice yaml.MapSlice) bool {

				fooIdx := slices.IndexFunc(mapSlice, func(item yaml.MapItem) bool { return item.Key == "foo" })
				bingIdx := slices.IndexFunc(mapSlice, func(item yaml.MapItem) bool { return item.Key == "bing" })
				bangIdx := slices.IndexFunc(mapSlice[bingIdx].Value.(yaml.MapSlice), func(item yaml.MapItem) bool { return item.Key == "bang" })
				boomIdx := slices.IndexFunc(mapSlice[bingIdx].Value.(yaml.MapSlice)[bangIdx].Value.(yaml.MapSlice), func(item yaml.MapItem) bool { return item.Key == "boom" })

				return fooIdx != -1 &&
					len(mapSlice) == 2 &&
					mapSlice[fooIdx].Value == "Bar" &&
					mapSlice[bingIdx].Value.(yaml.MapSlice)[bangIdx].Value.(yaml.MapSlice)[boomIdx].Value == "Bop"
			},
			frenchOk: func(mapSlice yaml.MapSlice) bool {
				fooIdx := slices.IndexFunc(mapSlice, func(item yaml.MapItem) bool { return item.Key == "foo" })
				bingIdx := slices.IndexFunc(mapSlice, func(item yaml.MapItem) bool { return item.Key == "bing" })
				bangIdx := slices.IndexFunc(mapSlice[bingIdx].Value.(yaml.MapSlice), func(item yaml.MapItem) bool { return item.Key == "bang" })
				boomIdx := slices.IndexFunc(mapSlice[bingIdx].Value.(yaml.MapSlice)[bangIdx].Value.(yaml.MapSlice), func(item yaml.MapItem) bool { return item.Key == "boom" })

				return fooIdx != -1 &&
					len(mapSlice) == 2 &&
					mapSlice[fooIdx].Value == "La Bar" &&
					mapSlice[bingIdx].Value.(yaml.MapSlice)[bangIdx].Value.(yaml.MapSlice)[boomIdx].Value == "La Bop"
			},
		},
		{
			translations: map[string]*Translation{
				"bing.bang.bap": {
					Key:         "bing.bang.bap",
					SourceValue: "Zip",
					Translations: map[string]string{
						"french": "La Zip",
					},
					Languages: []string{"french"},
				},
				"bing.bang.boom": {
					Key:         "bing.bang.boom",
					SourceValue: "Bop",
					Translations: map[string]string{
						"french": "La Bop",
					},
					Languages: []string{"french"},
				},
			},
			englishOk: func(mapSlice yaml.MapSlice) bool {
				return len(mapSlice) == 1
			},
			frenchOk: func(mapSlice yaml.MapSlice) bool { return true },
		},
	}

	for i, input := range inputs {
		if !input.englishOk(translationsToMapSlice("english", input.translations, true)) {
			t.Fatalf("Failed to transform translations in to map of english for input %d", i)
		}
		if !input.frenchOk(translationsToMapSlice("french", input.translations, false)) {
			t.Fatalf("Failed to transform translations in to map of french for input %d", i)
		}
	}

}

func TestToRelativePathPartsSimple(t *testing.T) {

	basePath := []string{
		"foo",
		"bar",
	}

	path := []string{
		"foo",
		"bar",
		"baz",
	}

	result := toRelativePathParts(basePath, path)

	if len(result) != 1 || result[0] != "baz" {
		t.Fatalf("Expected resulting path to be 'baz'")
	}
}

func TestToRelativePathPartsLong(t *testing.T) {

	basePath := []string{
		"foo",
		"bar",
	}

	path := []string{
		"foo",
		"bar",
		"baz",
		"biz",
		"bob",
	}

	result := toRelativePathParts(basePath, path)

	if len(result) != 3 || result[0] != "baz" || result[1] != "biz" || result[2] != "bob" {
		t.Fatalf("Expected resulting path to be 'baz.biz.bob'")
	}
}

func TestToRelativePathPartsBaseLongerThanPath(t *testing.T) {

	basePath := []string{
		"foo",
		"bar",
		"baz",
		"biz",
		"bob",
	}

	path := []string{
		"foo",
		"bar",
	}

	result := toRelativePathParts(basePath, path)

	if len(result) != 0 {
		t.Fatalf("Expected resulting path to be empty")
	}
}

func TestToRelativePathPartsBaseDivergentPath(t *testing.T) {

	basePath := []string{
		"blip",
	}

	path := []string{
		"foo",
		"bar",
	}

	result := toRelativePathParts(basePath, path)

	if len(result) != 0 {
		t.Fatalf("Expected resulting path to be empty")
	}
}
