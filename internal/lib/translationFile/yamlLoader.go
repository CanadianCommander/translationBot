package translationFile

import (
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type YamlLoader struct {
}

func (y *YamlLoader) Load(
	sourceLanguage string,
	translationLanguages []string,
	language string,
	file string,
	translations map[string]*Translation) (map[string]*Translation, error) {

	yamlRaw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var yamlData yaml.MapSlice = make([]yaml.MapItem, 0)
	err = yaml.Unmarshal(yamlRaw, &yamlData)
	if err != nil {
		return nil, err
	}

	extractTranslationsMapSlice(
		sourceLanguage,
		translationLanguages,
		language,
		"",
		yamlData,
		translations)

	return translations, nil
}

func (y *YamlLoader) CanLoad(filePath string) bool {
	return path.Ext(filePath) == ".yaml" || path.Ext(filePath) == ".yml"
}
