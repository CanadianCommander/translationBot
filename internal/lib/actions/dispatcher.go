package actions

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
	"time"
)

//==========================================================================
// Public
//==========================================================================

// Dispatch incoming action to appropriate handler
func Dispatch(interactionCallback *slack.InteractionCallback) error {

	for _, block := range interactionCallback.ActionCallback.BlockActions {
		go dispatchAction(interactionCallback, block)
	}

	return nil
}

//==========================================================================
// Private
//==========================================================================

// dispatchAction dispatches an action to the appropriate handler
func dispatchAction(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Errorf("Action handler %s panicked with error %s", block.ActionID, r)
			errorNotification(interactionCallback, "I'm freaking out man! "+slackutil.GetRandomEmoji())
		}
	}()

	var err error = nil

	log.Logger.Infof("Handling action %s", block.ActionID)
	startTime := time.Now()
	switch block.ActionID {
	case routes.ActionListProjects:
		err = listProjects(interactionCallback)
	case routes.ActionListMissingTranslations:
		err = ListMissingTranslations(interactionCallback)
	case routes.ActionUpdateTranslations:
		err = UpdateTranslations(interactionCallback, block)
	}
	log.Logger.Infof("%s handler completed in %dms", block.ActionID, time.Now().Sub(startTime).Milliseconds())

	if err != nil {
		log.Logger.Errorf("Unexpected error while handling slack action %s. %s", block.ActionID, err)
		errorNotification(interactionCallback, "Please contact @Benjamin Benetti")
	}
}

// errorNotification shows a nice little error message to the user.
func errorNotification(interactionCallback *slack.InteractionCallback, msg string) {
	slackutil.DeleteMessageByResponseUrl(interactionCallback.Channel.ID, interactionCallback.ResponseURL)
	slackutil.Api.PostEphemeral(
		interactionCallback.Channel.ID,
		interactionCallback.User.ID,
		slack.MsgOptionBlocks(
			slack.NewHeaderBlock(slackutil.NewTextBlock("That's not right!?!?!")),
			slack.NewSectionBlock(
				slackutil.NewTextBlock(fmt.Sprintf("I'm sorry I can't complete your request right now.\n %s", msg)),
				nil,
				nil)))
}
