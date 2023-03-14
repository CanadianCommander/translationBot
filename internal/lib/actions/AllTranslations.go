package actions

import (
	"errors"
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/configuration"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/ui"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ListAllTranslations produce all translations in the project specified as a CSV.
func ListAllTranslations(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) error {
	projectName, _, err := readProxyResponse(block.Value)
	if err != nil {
		return err
	}

	project := configuration.Get().GetProject(projectName)
	defer project.Unlock()

	action := getBlockActionById(routes.ActionListAllTranslations, interactionCallback)
	if action == nil {
		return errors.New("could not find action matching id " + routes.ActionListAllTranslations)
	}

	err = slackutil.DeleteMessageByResponseUrl(interactionCallback.Channel.ID, interactionCallback.ResponseURL)
	if err != nil {
		return err
	}

	response := ui.CsvDownload(project.Name, fmt.Sprintf("Click hear to download the complete translation set for %s", project.Name), "translations/csv")
	_, _, err = slackutil.Api.PostMessage(interactionCallback.Channel.ID, slackutil.SlackMessageToMsgOption(response))
	if err != nil {
		return err
	}

	return nil
}
