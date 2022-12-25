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
	blocks := []slack.Block{
		slack.NewHeaderBlock(slackutil.NewTextBlock("Missing translations")),
		slack.NewSectionBlock(
			slackutil.NewTextBlock("The following sections detail the missing translations for each language."+
				" If a language does not show up then it has no missing translations."+
				" Each line represents a single missing translation. It indicates the key and the english value."),
			nil,
			nil),
		slack.NewDividerBlock(),
	}
	blocks = append(blocks, getMissingTranslationOutputForLangs(missingTranslations)...)

	return slack.NewBlockMessage(blocks...)
}

//==========================================================================
// Private
//==========================================================================

// getMissingTranslationOutputForLangs gets missing translation string output messages by language
func getMissingTranslationOutputForLangs(missingTranslations []translation.Translation) []slack.Block {
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

	var outputBlocks []slack.Block
	for lang, out := range missingTranslationOutput {
		outputBlocks = append(outputBlocks, buildLanguageMissingTranslationBlocks(lang, out)...)
	}
	return outputBlocks
}

func buildLanguageMissingTranslationBlocks(language string, missingTranslations string) []slack.Block {
	return []slack.Block{
		slack.NewHeaderBlock(slackutil.NewTextBlock(language)),
		slack.NewSectionBlock(
			slackutil.NewMarkdownTextBlock(fmt.Sprintf("```%s\n```", missingTranslations)),
			nil,
			nil),
	}
}

// formatMissingTranslationLine formats a single line of missing translation output
// #### params
// translation - the translation to output the missing line for
func formatMissingTranslationLine(translation *translation.Translation) string {
	return fmt.Sprintf("\n%s \"%s\"", translation.Key, translation.SourceValue)
}
