package git

type LanguagePack struct {
	Name             string            `yaml:"name"`
	TranslationFiles map[string]string `yaml:"translationFiles"`
}
