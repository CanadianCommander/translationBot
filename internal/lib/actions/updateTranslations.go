package actions

import (
	"errors"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translationMapping"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

func UpdateTranslations(interactionCallback *slack.InteractionCallback, block *slack.BlockAction) error {
	log.Logger.Infof("Applying translation update using translation file %s", block.Value)

	loader, err := translationMapping.GetMapperForSlackFile(block.Value)
	if err != nil {
		return errors.New("unexpected error when searching for file loader")
	}

	slackFileReader, err := slackutil.DownloadSlackFileById(block.Value)
	if err != nil {
		return err
	}
	defer slackFileReader.Close()

	mappings, err := loader.Load(slackFileReader)
	if err != nil {
		return err
	}

	log.Logger.Infof("%+v", mappings)

	if err != nil {
		return err
	}

	return nil
}
