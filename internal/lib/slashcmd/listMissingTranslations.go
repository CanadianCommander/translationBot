package slashcmd

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// listMissingTranslations lists missing translations to the user.
func listMissingTranslations(slashCommand *slack.SlashCommand, args []string) error {

	var project *git.Project

	if len(args) > 1 {
		project = configuration.Get().GetProject(args[1])
		if project == nil {
			_, _, err := slackutil.Api.PostMessage(
				slashCommand.ChannelID,
				slackutil.SlackMessageToMsgOption(ui.ErrorNotification(fmt.Sprintf("Project %s doesn't exist", args[1]))))
			if err != nil {
				return err
			}
			return nil
		}
	} else {
		project = configuration.Get().GetDefaultProject()
	}
	defer project.Unlock()

	_, loaderId, err := slackutil.Api.PostMessage(
		slashCommand.ChannelID,
		slackutil.SlackMessageToMsgOption(ui.LoadingIndicator("Re-spooling the flux capacitors...")))
	if err != nil {
		return err
	}

	missingTranslations, err := translation.FindMissingTranslations(project)
	if err != nil {
		return err
	}

	_, _, err = slackutil.Api.DeleteMessage(slashCommand.ChannelID, loaderId)
	if err != nil {
		return err
	}

	response := ui.MissingTranslations(missingTranslations)
	_, _, err = slackutil.Api.PostMessage(slashCommand.ChannelID, slackutil.SlackMessageToMsgOption(response))
	if err != nil {
		return err
	}

	return nil
}
