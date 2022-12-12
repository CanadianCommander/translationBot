package main

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"os"
)
import "github.com/CanadianCommander/translationBot/internal/api"

func main() {
	log.InitializeLogging()
	fmt.Println(bootUpMessage())

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
