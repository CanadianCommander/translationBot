package configuration

import (
	"gopkg.in/yaml.v2"
	"os"
)

const configFilePath = "./config/config.yaml"

//==========================================================================
// Private
//==========================================================================

// loadSettings from the configuration file.
func loadSettings() *Settings {
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

		return &settings
	}
	panic("Config file doesnt exist. Freak out! Please make sure it exists and is at path ./config/config.yaml")
}
