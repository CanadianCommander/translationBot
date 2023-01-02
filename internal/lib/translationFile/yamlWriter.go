package translationFile

import (
	"errors"
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type YamlWriter struct {
}

func (yamlWriter *YamlWriter) Write(filePath string, lang string, sourceLanguage string, translations map[string]Translation) error {
	if !yamlWriter.CanWrite(filePath) {
		return errors.New(fmt.Sprintf("translationFile.YamlWriter does not support this type of file %s", filePath))
	}

	rawTranslations := translationsToMap(lang, translations, lang == sourceLanguage)

	yamlStr, err := yaml.Marshal(rawTranslations)
	if err != nil {
		log.Logger.Error("Error Marshaling translation YAML during translation file write, ", err)
		return err
	}

	file, err := os.OpenFile(filePath, os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(yamlStr)
	if err != nil {
		log.Logger.Error("Failed to write translation updates to YAML file", err)
		return err
	}

	return nil
}

func (yamlWriter *YamlWriter) CanWrite(filePath string) bool {
	return path.Ext(filePath) == ".yaml" || path.Ext(filePath) == ".yml"
}
