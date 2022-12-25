package slashcmd

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

type handler struct {
	Cmd     string
	Handler func(command slack.SlashCommand) slack.Message
}

// list of slash command handlers!
var handlers = []handler{
	{
		Cmd:     "",
		Handler: func(command slack.SlashCommand) slack.Message { return ui.Index() },
	},
}

// DispatchCommand to the appropriate handler and return the response from that handler
func DispatchCommand(slashCommand slack.SlashCommand) slack.Message {
	log.Logger.Infof("Processing Slash Command %s", slashCommand.Text)

	switch slashCommand.Text {
	case "":
		return ui.Index()
	default:
		return ui.NotFound()
	}
}
