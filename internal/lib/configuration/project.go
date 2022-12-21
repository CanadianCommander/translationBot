package configuration

type Project struct {
	Name        string
	Url         string
	GithubToken string

	// translation file mappings.
	// LANG: PATH
	// Ex english: ./foo/bar/bang.yaml
	TranslationFiles map[string]string
}
