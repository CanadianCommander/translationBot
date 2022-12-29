package actions

import (
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

func UpdateTranslations(action *slack.InteractionCallback, block *slack.BlockAction) error {
	log.Logger.Infof("Applying translation update using translation file %s", block.Value)

	//loader, err := translationMapping.GetMapperForSlackFile(block.Value)
	//if err != nil {
	//	return errors.New("unexpected error when searching for file loader")
	//}

	return nil
}
