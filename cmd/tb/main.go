package main

import (
	"fmt"
	"os"
)
import "github.com/CanadianCommander/translationBot/internal/api"

func main() {
	fmt.Println("Booting up Translation Bot")

	err := api.BuildV1Api().Run(":8080")
	if err != nil {
		fmt.Println("Failed to startup TranslationBot :(")
		os.Exit(1)
	}
}
