package translationFile

import (
	"bytes"
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"os"
	"os/exec"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type TsLoader struct {
}

func (t TsLoader) Load(
	sourceLanguage string,
	translationLanguages []string,
	language string,
	file string,
	translations map[string]Translation) (map[string]Translation, error) {

	cmdOutput := bytes.Buffer{}
	cmd := exec.Command("yarn", "--silent", "run", "toJson", file)
	cwd, _ := os.Getwd()
	cmd.Dir = path.Join(cwd, "cmd/tsJson/")
	cmd.Stdout = &cmdOutput
	if err := cmd.Run(); err != nil {
		log.Logger.Error("Error while running ", cmd.String())
		return nil, err
	}

	jsonData := make(map[string]interface{})
	if err := json.Unmarshal(cmdOutput.Bytes(), &jsonData); err != nil {
		return nil, err
	}

	extractTranslations(sourceLanguage, translationLanguages, language, "", jsonData, translations)
	return translations, nil
}

func (t TsLoader) CanLoad(filePath string) bool {
	return path.Ext(filePath) == ".ts"
}
