package tui

import (
	"github.com/luizvilasboas/commit-hooks/internal/config"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

const (
	stepInitialConfirm = iota
	stepType
	stepScope
	stepSubject
	stepBody
	stepBreakingChangePrompt
	stepBreakingChangeDescription
	stepConfirm
)

type model struct {
	initialConfirmChoices      []string
	commitTypeChoices          []string
	scopeChoices               []string
	subjectInput               textinput.Model
	bodyInput                  textarea.Model
	breakingChangePrompt       []string
	breakingChangeInput        textarea.Model
	confirmChoices             []string
	step                       int
	initialMessage             string
	initialConfirmCursor       int
	commitTypeCursor           int
	scopeCursor                int
	selectedScopes             map[int]struct{}
	breakingChangePromptCursor int
	confirmCursor              int
	quitting                   bool
	finalOutput                string
	errorMessage               string
	config                     config.Config
}

func newModel(initialMsg string, config config.Config) model {
	subjectInput := textinput.New()
	subjectInput.Placeholder = "Add new login endpoint"
	subjectInput.CharLimit = 100
	subjectInput.Width = 70
	subjectInput.SetValue(initialMsg)

	bodyInput := textarea.New()
	bodyInput.Placeholder = "More detailed description of the changes (optional)..."
	bodyInput.SetWidth(70)
	bodyInput.SetHeight(5)

	breakingChangeInput := textarea.New()
	breakingChangeInput.Placeholder = "Describe the breaking change..."
	breakingChangeInput.SetWidth(70)
	breakingChangeInput.SetHeight(5)

	m := model{
		initialMessage:        initialMsg,
		initialConfirmChoices: []string{"Yes, enhance this message", "No, use the original message"},
		commitTypeChoices:     config.Types,
		scopeChoices:          config.Scopes,
		selectedScopes:        make(map[int]struct{}),
		subjectInput:          subjectInput,
		bodyInput:             bodyInput,
		breakingChangePrompt:  []string{"No", "Yes"},
		breakingChangeInput:   breakingChangeInput,
		confirmChoices:        []string{"Confirm and Create Commit", "Start Over"},
		config:                config,
	}

	if initialMsg != "" {
		m.step = stepInitialConfirm
	} else {
		m.step = stepType
	}

	return m
}
