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

	if t.shouldInjectTemporaryPackage(file) {
		log.Logger.Info("Injecting temporary package to switch code to ESM module...")
		temporaryPackage, err := t.writeTemporaryPackageFile(file)
		if err != nil {
			return nil, err
		}
		defer t.cleanupTemporaryPackageFile(temporaryPackage)
	}

	cmdOutput := bytes.Buffer{}
	cmdErrorOut := bytes.Buffer{}
	cmd := exec.Command("yarn", "--silent", "run", "toJson", file)
	cwd, _ := os.Getwd()
	cmd.Dir = path.Join(cwd, "cmd/tsJson/")
	cmd.Stdout = &cmdOutput
	cmd.Stderr = &cmdErrorOut
	if err := cmd.Run(); err != nil {
		log.Logger.Error("Error while running ", cmd.String())
		log.Logger.Error(string(cmdErrorOut.Bytes()))
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
	return path.Ext(filePath) == ".ts" || path.Ext(filePath) == ".js"
}

//==========================================================================
// Private
//==========================================================================

// shouldInjectTemporaryPackage determines if a temporary package injection is required or not
func (t TsLoader) shouldInjectTemporaryPackage(filePath string) bool {
	_, err := os.Stat(path.Join(path.Dir(filePath), "package.json"))
	return err != nil
}

// writeTemporaryPackageFile creates a temporary package.json file in the directory of the given filePath
// This file will switch the JavaScript in that folder to module packaging mode.
func (t TsLoader) writeTemporaryPackageFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(path.Join(path.Dir(filePath), "package.json"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	_, err = file.Write([]byte("{\n  \"type\": \"module\"\n}\n"))
	return file, err
}

// cleanupTemporaryPackageFile deletes the temporary package file and closes its stream.
func (t TsLoader) cleanupTemporaryPackageFile(temporaryPackage *os.File) {

	err := temporaryPackage.Close()
	if err != nil {
		log.Logger.Error("Error cleaning up tmp package file", err)
	}

	err = os.Remove(temporaryPackage.Name())
	if err != nil {
		log.Logger.Error("Error cleaning up tmp package file", err)
	}
}
