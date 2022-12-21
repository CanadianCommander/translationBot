package configuration

import (
	"os"
	"sync"
)

var config *Settings = nil
var configMutex = sync.Mutex{}

func Get() *Settings {
	if config == nil {
		loadSettings()
	}

	return config
}

func loadSettings() {
	configMutex.Lock()
	defer configMutex.Unlock()

	slackSigningSecret, _ := os.LookupEnv("SLACK_SIGNING_KEY")

	config = &Settings{
		SlackSigningSecret: slackSigningSecret,
	}
}
