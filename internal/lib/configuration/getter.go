package configuration

import (
	"sync"
)

var config *Settings = nil
var configMutex = sync.Mutex{}

func Get() *Settings {
	if config == nil {
		configMutex.Lock()
		defer configMutex.Unlock()

		config = loadSettings()
	}
	return config
}
