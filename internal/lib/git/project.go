package git

import (
	"golang.org/x/exp/slices"
	"os"
	"path"
	"sync"
)

//==========================================================================
// Public
//==========================================================================

type Project struct {
	Name        string
	BaseDir     string
	Url         string
	Branch      string
	GitUsername string `yaml:"gitUsername"`
	GitPassword string `yaml:"gitPassword"`
	// The given command will run (in bash) before TranslationBot commits any changes
	PreCommitAction string `yaml:"preCommitAction"`

	// language that is the source of all other translations. aka "english"
	SourceLanguage string `yaml:"sourceLanguage"`
	// if true TranslationBot will update the source language files in addition to the translation files.
	UpdateSourceFile bool `yaml:"updateSourceFiles"`

	// translation "packs". array of one or more language packs.
	// a pack is composed of translation file mappings.
	// LANG: PATH
	// Ex
	// packs:
	//   - name: foo
	//     translationFiles:
	//       english: ./foo/bar/bang.yaml
	// 	     french: ./foo/bar/la-bang.yaml
	Packs []*LanguagePack `yaml:"packs"`

	// used to make sure only one goroute operates on a project at a time. The project folder is a shared resource.
	IsLocked    bool
	projectLock sync.Mutex
}

// ProjectRelativePathToAbsolute converts a path that is relative to the projects root to an absolute path on the system.
func (project *Project) ProjectRelativePathToAbsolute(filePath string) string {
	cwd, _ := os.Getwd()
	return path.Join(cwd, project.FilePath(), filePath)
}

//==========================================================================
// sync.Locker implementation
//==========================================================================

func (project *Project) Lock() {
	project.projectLock.Lock()
	project.IsLocked = true
}

func (project *Project) Unlock() {
	project.projectLock.Unlock()
	project.IsLocked = false
}

//==========================================================================
// Getters
//==========================================================================

func (project *Project) FilePath() string {
	return project.BaseDir
}

// TranslationFileCount counts the number of translation files in this project
func (project *Project) TranslationFileCount() int {

	numFiles := 0
	for _, pack := range project.Packs {
		numFiles += len(pack.TranslationFiles)
	}

	return numFiles
}

// TranslationLanguages returns the set of translation languages (all languages minus the source lang) for this project.
func (project *Project) TranslationLanguages() []string {
	var transLangs []string

	for _, pack := range project.Packs {
		for lang, _ := range pack.TranslationFiles {
			if lang != project.SourceLanguage && !slices.Contains(transLangs, lang) {
				transLangs = append(transLangs, lang)
			}
		}
	}

	return transLangs
}
