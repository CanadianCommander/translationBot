package configuration

type Settings struct {
	SlackSigningSecret string `yaml:"slackSigningSecret"`

	Projects []Project
}
