package ui

import (
	"github.com/CanadianCommander/translationBot/internal/lib/git"
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
func ProjectSelectList(projects []*git.Project, actions []string, values []string) slack.Message {

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

	blocks := []slack.Block{
		slack.NewHeaderBlock(slackutil.NewTextBlock("Which project do you want to run on?")),
		slack.NewDividerBlock(),
	}

	return slack.NewBlockMessage(append(blocks, projectBlocks...)...)
}
