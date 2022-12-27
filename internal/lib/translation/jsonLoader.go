package translation

import (
	"encoding/json"
	"os"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type JsonLoader struct {
}

func (j *JsonLoader) Load(
	sourceLanguage string,
	translationLanguages []string,
	language string,
	file string,
	translations map[string]Translation,
) (map[string]Translation, error) {
	rawJson, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	jsonData := make(map[string]interface{})
	err = json.Unmarshal(rawJson, &jsonData)
	if err != nil {
		return nil, err
	}

	extractTranslations(sourceLanguage, translationLanguages, language, "", jsonData, translations)
	return translations, nil
}

func (j *JsonLoader) CanLoad(filePath string) bool {
	return path.Ext(filePath) == ".json"
}
