package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/CanadianCommander/translationBot/internal/lib/translationFile"
	"github.com/slack-go/slack"
	"sort"
	"strings"
)

//==========================================================================
// Public
//==========================================================================

// MissingTranslations outputs a slack message that allows users to visualize the missing translations.
// #### params
// missingTranslations - a list of translations missing values for one or more languages
// projectName - the name of the project to which the missing translations pertain
func MissingTranslations(missingTranslations []translationFile.Translation, projectName string) slack.Message {
	blocks := []slack.Block{
		slack.NewHeaderBlock(slackutil.NewTextBlock(fmt.Sprintf("%s Missing translations", projectName))),
		slack.NewContextBlock(
			"instructions",
			slackutil.NewTextBlock("The following sections detail the missing translations for each language."+
				" If a language does not show up then it has no missing translations."+
				" Each line represents a single missing translation. It indicates the translation key and the english string."),
		),
	}
	if len(missingTranslations) > 0 {
		blocks = append(blocks, getMissingTranslationOutputForLangs(missingTranslations)...)
	} else {
		blocks = append(
			blocks,
			slack.NewSectionBlock(
				slackutil.NewTextBlock("All translations up to date! :tada:"),
				nil,
				nil),
		)
	}

	return slack.NewBlockMessage(blocks...)
}

//==========================================================================
// Private
//==========================================================================

// getMissingTranslationOutputForLangs gets missing translation string output messages by language
func getMissingTranslationOutputForLangs(missingTranslations []translationFile.Translation) []slack.Block {
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

	var textBlocks []slack.Block
	var currTextBlock *slack.TextBlockObject
	lines := strings.Split(strings.Trim(missingTranslations, "\n"), "\n")

	sort.Strings(lines)
	for idx, line := range lines {
		formattedLine := fmt.Sprintf("- %s\n", line)

		if currTextBlock == nil || len(currTextBlock.Text)+len(formattedLine)+3 >= slackutil.TextBlockCharacterLimit {
			if currTextBlock != nil {
				currTextBlock.Text += "```"
			}

			currTextBlock = slackutil.NewMarkdownTextBlock(fmt.Sprintf("```%s", formattedLine))
			textBlocks = append(textBlocks, slack.NewSectionBlock(currTextBlock, nil, nil))
		} else {
			currTextBlock.Text += formattedLine
			if idx == len(lines)-1 {
				currTextBlock.Text += "```"
			}
		}
	}

	outputBlocks := []slack.Block{
		slack.NewDividerBlock(),
		slack.NewHeaderBlock(slackutil.NewTextBlock(strings.ToUpper(language[:1]) + language[1:])),
	}
	outputBlocks = append(outputBlocks, textBlocks...)

	return outputBlocks
}

// formatMissingTranslationLine formats a single line of missing translation output
// #### params
// translation - the translation to output the missing line for
func formatMissingTranslationLine(translation *translationFile.Translation) string {
	return fmt.Sprintf("\n%s \"%s\"", translation.Key, translation.SourceValue)
}
