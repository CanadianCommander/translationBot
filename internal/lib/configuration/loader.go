package configuration

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

const configFilePath = "./config/config.yaml"

//==========================================================================
// Private
//==========================================================================

const EnvironmentVersion = "VERSION"

// loadSettings from the configuration file.
func loadSettings() *Settings {
	log.Logger.Infof("Loading configuration file %s", configFilePath)

	fileInfo, err := os.Stat(configFilePath)
	if fileInfo != nil || err != nil {
		configYaml, err := os.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}

		var settings = Settings{}
		err = yaml.Unmarshal(configYaml, &settings)
		if err != nil {
			panic(err)
		}

		settings.Version = os.Getenv(EnvironmentVersion)
		if settings.Version == "" {
			settings.Version = "unknown"
		}

		log.Logger.Info("Configuration loaded")
		return postProcessProjectSettings(&settings)
	}
	panic("Config file doesnt exist. Freak out! Please make sure it exists and is at path ./config/config.yaml")
}

// postProcessProjectSettings applies post-processing to project settings
func postProcessProjectSettings(settings *Settings) *Settings {
	for key, val := range settings.Projects {
		val.Name = key
		val.BaseDir = path.Join(git.RepoStorageLocation, val.Name)

		for _, pack := range val.Packs {
			pack.Project = val
		}
	}

	return settings
}
