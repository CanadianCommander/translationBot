package translationFile

import (
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
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
	pack *git.LanguagePack,
	translations map[string]*Translation,
) (map[string]*Translation, error) {
	rawJson, err := os.ReadFile(pack.Project.ProjectRelativePathToAbsolute(pack.TranslationFiles[language]))
	if err != nil {
		return nil, err
	}

	jsonData := make(map[string]interface{})
	err = json.Unmarshal(rawJson, &jsonData)
	if err != nil {
		return nil, err
	}

	extractTranslations(pack, sourceLanguage, translationLanguages, language, "", jsonData, translations)
	return translations, nil
}

func (j *JsonLoader) CanLoad(filePath string) bool {
	return path.Ext(filePath) == ".json"
}
