package git

import (
	"os"
	"path"
)

//==========================================================================
// Public
//==========================================================================

type Project struct {
	Name        string
	Url         string
	Branch      string
	GitUsername string `yaml:"gitUsername"`
	GitPassword string `yaml:"gitPassword"`

	// language that is the source of all other translations. aka "english"
	SourceLanguage string `yaml:"sourceLanguage"`
	// translation file mappings.
	// LANG: PATH
	// Ex english: ./foo/bar/bang.yaml
	TranslationFiles map[string]string `yaml:"translationFiles"`
}

// ProjectRelativePathToAbsolute converts a path that is relative to the projects root to an absolute path on the system.
func (project *Project) ProjectRelativePathToAbsolute(filePath string) string {
	cwd, _ := os.Getwd()
	return path.Join(cwd, project.FilePath(), filePath)
}

//==========================================================================
// Getters
//==========================================================================

func (project *Project) FilePath() string {
	return path.Join(RepoStorageLocation, project.Name)
}

// TranslationLanguages returns the set of translation languages (all languages minus the source lang) for this project.
func (project *Project) TranslationLanguages() []string {
	var transLangs []string

	for lang, _ := range project.TranslationFiles {
		if lang != project.SourceLanguage {
			transLangs = append(transLangs, lang)
		}
	}

	return transLangs
}
