package translationMapping

import (
	"bytes"
	_ "embed"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"golang.org/x/exp/slices"
	"testing"
)

//go:embed test/testTranslationFile.csv
var testTranslationFile string

type loaderTestInput struct {
	FileData string
	resultOk func(tm []translationFile.Translation, err error) bool
}

func TestCsvTaraMappingLoader_Load(t *testing.T) {

	csvTara := &CsvTaraMappingLoader{}

	testInputs := []*loaderTestInput{
		{
			FileData: "ENGLISH,FRENCH \nGo to the beach,La go to the beach",
			resultOk: func(tm []translationFile.Translation, err error) bool {
				return len(tm) == 1 &&
					tm[0].SourceValue == "Go to the beach" &&
					tm[0].Translations["french"] == "La go to the beach" &&
					err == nil
			},
		},
		{
			FileData: "ENGLISH,FRENCH,WELSH\nI like beer,I like wine,Dw i'n hoffi cwrw",
			resultOk: func(tm []translationFile.Translation, err error) bool {
				return len(tm) == 1 &&
					tm[0].SourceValue == "I like beer" &&
					tm[0].Translations["french"] == "I like wine" &&
					tm[0].Translations["welsh"] == "Dw i'n hoffi cwrw" &&
					err == nil
			},
		},
		{
			FileData: "ENGLISH,FRENCH\n\"Commas are cool,,,,,, see! so cool,,,,\",\"La commas are cool,,,,, ,,\"",
			resultOk: func(tm []translationFile.Translation, err error) bool {
				return len(tm) == 1 &&
					tm[0].SourceValue == "Commas are cool,,,,,, see! so cool,,,," &&
					tm[0].Translations["french"] == "La commas are cool,,,,, ,," &&
					err == nil
			},
		},
		{
			FileData: testTranslationFile,
			resultOk: func(tm []translationFile.Translation, err error) bool {
				translationIndex := slices.IndexFunc(tm, func(t translationFile.Translation) bool { return t.SourceValue == "Add New Card" })

				return len(tm) == 826 &&
					tm[translationIndex].Translations["french"] == "Ajoutez une nouvelle carte" &&
					err == nil
			},
		},
	}

	for i, input := range testInputs {
		if !input.resultOk(csvTara.Load(bytes.NewBuffer([]byte(input.FileData)))) {
			t.Fatalf("CSVTara failed to load test input file %d", i)
		}
	}
}

func TestCsvTaraMappingLoader_IsFileSupported_Supported(t *testing.T) {

	csvTara := &CsvTaraMappingLoader{}

	testFiles := []string{
		"ENGLISH,FRENCH\nfoo,bar",
		"ENGLISH, FRENCH\nbob,ross",
		"ENGLISH      , FRENCH\nbob,ross",
		"         ENGLISH      , FRENCH    \nbob,ross, frank, drebin",
	}

	for _, testFile := range testFiles {
		if !csvTara.IsFileSupported([]byte(testFile)) {
			t.Fatalf("CSVTara1.0 does not support file it should support %s", testFile)
		}
	}
}

func TestCsvTaraMappingLoader_IsFileSupported_NotSupported(t *testing.T) {

	csvTara := &CsvTaraMappingLoader{}

	testFiles := []string{
		"FLORP,BING\nfoo,bar",
		">>>>>ENGLISH,\\\\ FRENCH\nbob,ross",
		"This shiz isn't even a csv",
		"foo, bar,\n bang, fiz",
	}

	for _, testFile := range testFiles {
		if csvTara.IsFileSupported([]byte(testFile)) {
			t.Fatalf("CSVTara1.0 claimes to support file it shouldn't support %s", testFile)
		}
	}
}

func TestCsvTaraMappingLoader_IsMimeSupported_Supported(t *testing.T) {
	csvTara := &CsvTaraMappingLoader{}

	if !csvTara.IsMimeSupported("text/csv") {
		t.Fatalf("CSVTara1.0 indicated it did not support text/csv which is incorrect")
	}
}

func TestCsvTaraMappingLoader_IsMimeSupported_NotSupported(t *testing.T) {
	csvTara := &CsvTaraMappingLoader{}

	if csvTara.IsMimeSupported("image/csv") {
		t.Fatalf("CSVTara1.0 indicated it does support image/jpg which is incorrect")
	}
}
