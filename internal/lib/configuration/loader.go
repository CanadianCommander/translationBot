package configuration

import "os"

func Load() Settings {
	slackSigningSecret, _ := os.LookupEnv("SLACK_SIGNING_KEY")

	return Settings{
		SlackSigningSecret: slackSigningSecret,
	}
}
