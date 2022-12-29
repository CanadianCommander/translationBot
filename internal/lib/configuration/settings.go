package configuration

import "github.com/CanadianCommander/translationBot/internal/lib/git"

//==========================================================================
// Public
//==========================================================================

type Settings struct {
	SlackClientSecret  string `yaml:"slackClientSecret"`
	SlackSigningSecret string `yaml:"slackSigningSecret"`

	// what project should translation bot operate on if no project is specified in user command
	DefaultProject string `yaml:"defaultProject"`
	// projects on which translation bot can operate
	Projects map[string]*git.Project
}

//==========================================================================
// Getters
//==========================================================================

// GetDefaultProject returns the default project as configured.
// Said project is LOCKED. The caller must unlock the project before discarding it.
func (s *Settings) GetDefaultProject() *git.Project {
	defaultProject := s.Projects[s.DefaultProject]

	if defaultProject != nil {
		defaultProject.Lock()
		return defaultProject
	}
	return defaultProject
}
