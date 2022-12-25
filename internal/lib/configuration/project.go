package configuration

type Project struct {
	Url         string
	Branch      string
	GithubToken string `yaml:"githubToken"`

	// translation file mappings.
	// LANG: PATH
	// Ex english: ./foo/bar/bang.yaml
	TranslationFiles map[string]string `yaml:"translationFiles"`
}
