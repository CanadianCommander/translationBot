package slackutil

import (
	"encoding/json"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// MultiProjectInputDto is the input dto for the multi project proxy action
type MultiProjectInputDto struct {
	ActionId string
	Value    string
	// should the user be given the option to go back to the index page?
	ShowBackToIndex bool
}

// NewMultiProjectButtonBlockElement creates a new slack button element that allows the user to select which project
// the action performed by the button should operate on
// #### params
// actionId - the action to perform
// value - the value to send to that action
// text - the button text
func NewMultiProjectButtonBlockElement(actionId string, value string, text string, showBackToIndex bool) *slack.ButtonBlockElement {
	proxyValue := MultiProjectInputDto{
		ActionId:        actionId,
		Value:           value,
		ShowBackToIndex: showBackToIndex,
	}
	proxyValueStr, err := json.Marshal(proxyValue)
	if err != nil {
		log.Logger.Errorf("Serialization error while building multi project button block element!")
		panic(err)
	}

	return slack.NewButtonBlockElement(
		routes.ActionProjectProxy,
		string(proxyValueStr),
		NewTextBlock(text))
}
