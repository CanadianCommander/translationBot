package translationFile

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type TsWriter struct {
}

func (tsWriter *TsWriter) Write(filePath string, lang string, sourceLanguage string, translations map[string]*Translation) error {
	if !tsWriter.CanWrite(filePath) {
		return errors.New(fmt.Sprintf("translationFile.TsWriter does not support this type of file %s", filePath))
	}

	rawTranslations := translationsToMapSlice(lang, translations, lang == sourceLanguage)

	yamlStr, err := yaml.Marshal(rawTranslations)
	if err != nil {
		log.Logger.Error("Error Marshaling translation JSON during translation file write, ", err)
		return err
	}

	// write to ts file
	inputJsonBuffer := bytes.NewBuffer(yamlStr)
	cmd := exec.Command("yarn", "--silent", "run", "fromYaml", filePath)
	cwd, _ := os.Getwd()
	cmd.Dir = path.Join(cwd, "cmd/tsYaml/")
	cmd.Stdin = inputJsonBuffer
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Logger.Error("Error while running ", cmd.String())
		log.Logger.Error(out)
		return err
	}

	return nil
}

func (tsWriter *TsWriter) CanWrite(filePath string) bool {
	return path.Ext(filePath) == ".ts" || path.Ext(filePath) == ".js"
}
