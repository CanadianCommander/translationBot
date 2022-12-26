package git

import "path"

//==========================================================================
// Public
//==========================================================================

type Project struct {
	Name        string
	Url         string
	Branch      string
	GitUsername string `yaml:"gitUsername"`
	GitPassword string `yaml:"gitPassword"`

	// translation file mappings.
	// LANG: PATH
	// Ex english: ./foo/bar/bang.yaml
	TranslationFiles map[string]string `yaml:"translationFiles"`
}

//==========================================================================
// Getters
//==========================================================================

func (project *Project) filePath() string {
	return path.Join(RepoStorageLocation, project.Name)
}