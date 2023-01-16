package configuration

import "github.com/CanadianCommander/translationBot/internal/lib/git"

//==========================================================================
// Public
//==========================================================================

type Settings struct {
	// translation bot version.
	Hostname           string `yaml:"hostname"`
	Version            string
	SlackClientSecret  string `yaml:"slackClientSecret"`
	SlackSigningSecret string `yaml:"slackSigningSecret"`

	// what project should translation bot operate on if no project is specified in user command
	DefaultProject string `yaml:"defaultProject"`
	// if true all 'disruptive' remote actions will be skipped. i.e. translation updates will not be pushed and no PR will be opened
	TestMode bool `yaml:"testMode"`
	// projects on which translation bot can operate
	Projects map[string]*git.Project
}

//==========================================================================
// Getters
//==========================================================================

// GetDefaultProject returns the default project as configured.
// Said project is LOCKED. The caller must unlock the project before discarding it.
func (s *Settings) GetDefaultProject() *git.Project {
	return s.GetProject(s.DefaultProject)
}

// GetProject gets a project by name
func (s *Settings) GetProject(projectName string) *git.Project {
	project := s.Projects[projectName]

	if project != nil {
		project.Lock()
		return project
	}
	return project
}
