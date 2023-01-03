package main

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"os"
)
import "github.com/CanadianCommander/translationBot/internal/api"

func main() {
	fmt.Println(bootUpMessage())

	log.InitializeLogging()
	slackutil.InitializeApi()

	if configuration.Get().TestMode {
		log.Logger.Infof("Running in test mode")
	}

	err := api.BuildV1Api().Run(":8080")
	if err != nil {
		log.Logger.Error("Failed to startup TranslationBot :(")
		os.Exit(1)
	}
}

func bootUpMessage() string {
	return `
 ______                  __     __  _             ___       __ 
/_  __/______ ____  ___ / /__ _/ /_(_)__  ___    / _ )___  / /_
 / / / __/ _ ` + ` / _ \(_-</ / _ ` + ` / __/ / _ \/ _ \  / _  / _ \/ __/
/_/ /_/  \_,_/_//_/___/_/\_,_/\__/_/\___/_//_/ /____/\___/\__/ 
==============================================================
`
}
