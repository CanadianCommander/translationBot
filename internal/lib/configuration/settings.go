package configuration

type Settings struct {
	SlackClientSecret  string `yaml:"slackClientSecret"`
	SlackSigningSecret string `yaml:"slackSigningSecret"`

	// what project should translation bot operate on if no project is specified in user command
	DefaultProject string `yaml:"defaultProject"`
	// projects on which translation bot can operate
	Projects map[string]Project
}
