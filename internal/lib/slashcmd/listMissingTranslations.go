package slashcmd

import (
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// listMissingTranslations lists missing translations to the user.
func listMissingTranslations(slashCommand *slack.SlashCommand) error {
	project := configuration.Get().GetDefaultProject()
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
