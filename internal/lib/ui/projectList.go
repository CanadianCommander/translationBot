package ui

import (
	"fmt"
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Private
//==========================================================================

// ProjectList renders a ui showing the list of projects supported by TranslationBot
// #### params
// projects - projects to list
func ProjectList(projects []*git.Project) slack.Message {

	projectList := "```\n"
	for _, project := range projects {
		projectList = fmt.Sprintf("%s - %s\n", projectList, project.Name)
	}
	projectList += "```"

	return slack.NewBlockMessage(
		slack.NewHeaderBlock(slackutil.NewTextBlock("Supported Projects")),
		slack.NewDividerBlock(),
		slack.NewContextBlock("blurb",
			slackutil.NewTextBlock("I can help you with any of the following projects")),
		slack.NewSectionBlock(
			slackutil.NewMarkdownTextBlock(projectList),
			nil,
			nil),
	)

}
