package slashcmd

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
	"golang.org/x/exp/maps"
	"strings"
	"time"
)

//==========================================================================
// Public
//==========================================================================

// DispatchCommand to the appropriate handler and return the response from that handler
func DispatchCommand(slashCommand slack.SlashCommand) {
	log.Logger.Infof("Processing Slash Command %s", slashCommand.Text)
	config := configuration.Get()

	defer func() {
		if r := recover(); r != nil {
			log.Logger.Errorf("Slash command handler [%s] panicked with error %s", slashCommand.Text, r)
			_ = simpleResponse(&slashCommand, ui.ErrorNotification("O the humanity!"))
		}
	}()

	startTime := time.Now()
	args := strings.Split(slashCommand.Text, " ")
	var err error
	switch args[0] {
	case "":
		err = simpleResponse(&slashCommand, ui.Index())
	case "projects":
		err = simpleResponse(&slashCommand, ui.ProjectList(maps.Values(config.Projects)))
	case "missing":
		err = listMissingTranslations(&slashCommand, args)
	case "all":
		err = allTranslations(&slashCommand, args)
	case "release", "notes":
		err = simplePublicResponse(&slashCommand, ui.ReleaseNotes(false))
	default:
		err = simpleResponse(&slashCommand, ui.NotFound())
	}

	if err != nil {
		log.Logger.Errorf("Unexpected error while handling slash command %s\n %s", slashCommand.Text, err)
		_ = simpleResponse(&slashCommand, ui.ErrorNotification("To the logs!"))
	}

	log.Logger.Infof("Slash command %s processed in %dms", slashCommand.Text, time.Now().Sub(startTime).Milliseconds())
}
