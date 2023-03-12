package translationMapping

import (
	"bytes"
	_ "embed"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"testing"
)

type csvKeyedTaraLoaderTestInput struct {
	FileData string
	resultOk func(tm []translationFile.Translation, err error) bool
}

func TestCsvKeyedTaraMappingLoader_Load(t *testing.T) {

	csvTara := &CsvKeyedTaraMappingLoader{}

	testInputs := []*csvKeyedTaraLoaderTestInput{
		{
			FileData: "KEY,ENGLISH,FRENCH\n foo.bar, florp, la florp \n foo.biz, english, french",
			resultOk: func(tm []translationFile.Translation, err error) bool {
				return len(tm) == 2 &&
					tm[0].SourceValue == "florp" &&
					tm[0].Key == "foo.bar" &&
					tm[0].Translations["french"] == "la florp" &&
					tm[1].SourceValue == "english" &&
					tm[1].Key == "foo.biz" &&
					err == nil
			},
		},
		{
			FileData: "KEY,ENGLISH,FRENCH\n foo.bar, florp, \n foo.biz, english, french",
			resultOk: func(tm []translationFile.Translation, err error) bool {
				return len(tm) == 2 &&
					tm[0].SourceValue == "florp" &&
					tm[0].Key == "foo.bar" &&
					tm[0].Translations["french"] == "" &&
					tm[1].SourceValue == "english" &&
					tm[1].Key == "foo.biz" &&
					err == nil
			},
		},
	}

	for i, input := range testInputs {
		translationMap, err := csvTara.Load(bytes.NewBuffer([]byte(input.FileData)))
		if !input.resultOk(translationMap, err) {
			if err != nil {
				t.Fatalf("CSVKeyedTara failed to load test input file %d. With error: %s", i, err.Error())
			} else {
				t.Fatalf("CSVKeyedTara failed to load test input file %d", i)
			}
		}
	}
}

func TestCsvKeyedTaraMappingLoader_IsFileSupported_Supported(t *testing.T) {

	csvTara := &CsvKeyedTaraMappingLoader{}

	testFiles := []string{
		"KEY, ENGLISH,FRENCH\nflip.flop, foo,bar",
		" KEY ,ENGLISH, FRENCH\nbip.bop, bob,ross",
		"    KEY, ENGLISH      , FRENCH\n    drip.drop, bob,ross",
		"KEY    ,        ENGLISH   , FRENCH    \n apple.sauce, bob,ross, frank, drebin",
	}

	for _, testFile := range testFiles {
		if !csvTara.IsFileSupported([]byte(testFile)) {
			t.Fatalf("CSVKeyedTara1.0 does not support file it should support %s", testFile)
		}
	}
}

func TestCsvKeyedTaraMappingLoader_IsFileSupported_NotSupported(t *testing.T) {

	csvTara := &CsvKeyedTaraMappingLoader{}

	testFiles := []string{
		"FLORP,BING\nfoo,bar",
		">>>>>ENGLISH,\\\\ FRENCH\nbob,ross",
		"ENGLISH,FRENCH\nbob,ross",
		"This shiz isn't even a csv",
		"foo, bar,\n bang, fiz",
	}

	for _, testFile := range testFiles {
		if csvTara.IsFileSupported([]byte(testFile)) {
			t.Fatalf("CSVTara1.0 claimes to support file it shouldn't support %s", testFile)
		}
	}
}

func TestCsvKeyedTaraMappingLoader_IsMimeSupported_Supported(t *testing.T) {
	csvTara := &CsvKeyedTaraMappingLoader{}

	if !csvTara.IsMimeSupported("text/csv") {
		t.Fatalf("CSVKeyedTara1.0 indicated it did not support text/csv which is incorrect")
	}
}

func TestCsvKeyedTaraMappingLoader_IsMimeSupported_NotSupported(t *testing.T) {
	csvTara := &CsvKeyedTaraMappingLoader{}

	if csvTara.IsMimeSupported("image/csv") {
		t.Fatalf("CSVKeyedTara1.0 indicated it does support image/jpg which is incorrect")
	}
}
