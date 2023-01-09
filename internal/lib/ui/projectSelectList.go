package ui

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
	"github.com/CanadianCommander/translationBot/internal/lib/routes"
	"github.com/CanadianCommander/translationBot/internal/lib/slackutil"
	"github.com/slack-go/slack"
)

//==========================================================================
// Public
//==========================================================================

// ProjectSelectList renders a ui showing the list of projects from which the user can select
// #### params
// projects - projects to list
// actions - a list of actions in exactly the same order as projects. This list maps projects to slack actions
// values - a list of values in exactly the same order as projects. This list maps values to projects for use in actions
// showBackToIndex - should the back to index button display
func ProjectSelectList(projects []*git.Project, actions []string, values []string, showBackToIndex bool) slack.Message {

	var projectBlocks []slack.Block
	for idx, project := range projects {
		projectBlocks = append(
			projectBlocks,
			slack.NewSectionBlock(
				slackutil.NewTextBlock(project.Name),
				nil,
				&slack.Accessory{
					ButtonElement: slack.NewButtonBlockElement(
						actions[idx],
						values[idx],
						slackutil.NewTextBlock("Select")),
				}))
	}

	var actionButtons []slack.BlockElement
	if showBackToIndex {
		actionButtons = append(
			actionButtons,
			slack.NewButtonBlockElement(routes.ActionIndex, "", slackutil.NewTextBlock("Back")))
	} else {
		actionButtons = append(
			actionButtons,
			slack.NewButtonBlockElement(routes.ActionCancel, "", slackutil.NewTextBlock("Cancel")))
	}

	blocks := []slack.Block{
		slack.NewHeaderBlock(slackutil.NewTextBlock("Which project do you want to run on?")),
		slack.NewDividerBlock(),
	}
	blocks = append(blocks, projectBlocks...)
	blocks = append(blocks,
		slack.NewActionBlock(
			"actions",
			actionButtons...))

	return slack.NewBlockMessage(blocks...)
}
