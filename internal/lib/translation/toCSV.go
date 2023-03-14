package translation

import (
	"encoding/csv"
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"sort"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

// ToCSV converts a slice of translations into a CSV formatted string.
// #### params
// translations - translations to convert
// sourceLang - the source language
// #### return
// string of the format CSV format
// KEY, SOURCE, LANG1, LANG2, LANG3
// KEY, S_VAL, L1_VAL, L2_VAL, L3_VAL...
// ....
func ToCSV(translations []*translationFile.Translation, sourceLang string) (string, error) {
	if len(translations) == 0 {
		return "", errors.New("translation list cannot be empty when converting to CSV")
	}

	// sort alphabetically by keypath
	sort.Slice(translations, func(i int, j int) bool {
		return strings.Compare(translations[i].Key, translations[j].Key) < 0
	})

	firstTranslation := translations[0]
	stringBuilder := strings.Builder{}
	csvWriter := csv.NewWriter(&stringBuilder)

	allLanguages := []string{sourceLang}
	translationLanguages := firstTranslation.Languages
	sort.Strings(translationLanguages)
	allLanguages = append(allLanguages, translationLanguages...)

	headers := []string{"KEY"}
	headers = append(headers, allLanguages...)
	for i, header := range headers {
		headers[i] = strings.ToUpper(header)
	}

	err := csvWriter.Write(headers)
	if err != nil {
		return "", err
	}

	for _, trans := range translations {

		outputs := []string{trans.Key}
		for _, lang := range allLanguages {
			if lang == sourceLang {
				outputs = append(outputs, trans.SourceValue)
			} else {
				outputs = append(outputs, trans.Translations[lang])
			}
		}

		err = csvWriter.Write(outputs)
		if err != nil {
			return "", err
		}
	}

	csvWriter.Flush()
	return stringBuilder.String(), nil
}
