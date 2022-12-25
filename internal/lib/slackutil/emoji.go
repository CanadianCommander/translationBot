package slackutil

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/log"
	"golang.org/x/exp/maps"
	"math/rand"
)

//==========================================================================
// Public
//==========================================================================

// GetRandomEmoji pulls a random emoji from the list of emojis in the slack workspace.
func GetRandomEmoji() string {
	emojis, err := Api.GetEmoji()
	if err != nil {
		log.Logger.Error("Error getting emoji list defaulting to :robot:. Error: ", err)
		return ":robot:"
	}
	return fmt.Sprintf(":%s:", maps.Keys(emojis)[rand.Intn(len(emojis))])
}
