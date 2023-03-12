package translationFile

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"os"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type JsonWriter struct {
}

func (jsonWriter *JsonWriter) Write(filePath string, lang string, sourceLanguage string, translations map[string]*Translation) error {
	if !jsonWriter.CanWrite(filePath) {
		return errors.New(fmt.Sprintf("translationFile.JsonWriter does not support this type of file %s", filePath))
	}

	rawTranslations := translationsToMap(lang, translations, lang == sourceLanguage)

	jsonStr, err := json.MarshalIndent(rawTranslations, "", "  ")
	if err != nil {
		log.Logger.Error("Error Marshaling translation JSON during translation file write, ", err)
		return err
	}

	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonStr)
	if err != nil {
		log.Logger.Error("Failed to write translation updates to JSON file", err)
		return err
	}

	return nil
}

func (jsonWriter *JsonWriter) CanWrite(filePath string) bool {
	return path.Ext(filePath) == ".json"
}
