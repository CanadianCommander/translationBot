package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translation"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// MissingTranslations outputs a slack message that allows users to visualize the missing translations.
// #### params
// missingTranslations - a list of translations missing values for one or more languages
func MissingTranslations(missingTranslations []translation.Translation) slack.Message {

	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slackutil.NewTextBlock("Missing translations")),
		slack.NewSectionBlock(
			slackutil.NewTextBlock("The following sections detail the missing translations for each language."+
				" If a language does not show up then it has no missing translations."+
				" Each line represents a single missing translation. It indicates the key and the english value."),
			nil,
			nil),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			slackutil.NewMarkdownTextBlock(getMissingTranslationOutputForLangs(missingTranslations)),
			nil,
			nil),
	)
}

//==========================================================================
// Private
//==========================================================================

// getMissingTranslationOutputForLangs gets missing translation string output messages by language
func getMissingTranslationOutputForLangs(missingTranslations []translation.Translation) string {
	// lang -> output text
	missingTranslationOutput := make(map[string]string)

	// group outputs by language.
	for _, translation := range missingTranslations {
		if translation.HasMissingTranslations() {
			for _, missingLang := range translation.MissingLanguages() {
				output, exists := missingTranslationOutput[missingLang]
				if !exists {
					output = ""
				}
				missingTranslationOutput[missingLang] = output + formatMissingTranslationLine(&translation)
			}
		}
	}

	output := ""
	for lang, out := range missingTranslationOutput {
		output += fmt.Sprintf("\n*%s*\n```%s```", lang, out)
	}
	return output
}

func formatMissingTranslationLine(translation *translation.Translation) string {
	return fmt.Sprintf("\n%s \"%s\"", translation.Key, translation.SourceValue)
}
