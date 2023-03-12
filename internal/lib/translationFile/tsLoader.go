package translationFile

import (
	"bytes"
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"io/ioutil"
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
	translations map[string]*Translation) (map[string]*Translation, error) {

	if t.shouldInjectTemporaryPackage(file) {
		log.Logger.Info("Creating temporary package environment to switch code to ESM module...")
		tempDir, newFile, err := t.writeTemporaryPackageEnvironment(file)
		if err != nil {
			return nil, err
		}
		file = newFile

		defer t.cleanupTemporaryPackageEnvironment(tempDir)
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

// writeTemporaryPackageEnvironment creates a temporary package.json file in a temporary directory
// This file will switch the JavaScript in that folder to module packaging mode.
// #### return
// 0 - the temporary directory path
// 1 - the modified file path that was passed in. It will not reference a temporary copy
// 2 - error
func (t TsLoader) writeTemporaryPackageEnvironment(filePath string) (string, string, error) {
	tempDir, err := ioutil.TempDir("/tmp/", "jsenv")
	if err != nil {
		return "", "", err
	}

	file, err := os.OpenFile(path.Join(tempDir, "package.json"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return "", "", err
	}
	defer file.Close()
	_, err = file.Write([]byte("{\n  \"type\": \"module\"\n}\n"))

	// copy translation file data in to temp dir
	translationData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", "", err
	}
	newFilePath := path.Join(tempDir, path.Base(filePath))
	err = ioutil.WriteFile(newFilePath, translationData, 0666)
	if err != nil {
		return "", "", err
	}

	return tempDir, newFilePath, err
}

// cleanupTemporaryPackageFile deletes the temporary package environment.
func (t TsLoader) cleanupTemporaryPackageEnvironment(temporaryPackageDir string) {
	err := os.RemoveAll(temporaryPackageDir)
	if err != nil {
		log.Logger.Errorf("Error cleaning up tmp package environment %s. - %s", temporaryPackageDir, err)
	}
}
