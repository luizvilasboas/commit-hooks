package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMessage = ""
		if msg.Type == tea.KeyCtrlC || msg.String() == "q" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.step {
		case stepInitialConfirm:
			switch msg.String() {
			case "up", "k", "down", "j":
				m.initialConfirmCursor = (m.initialConfirmCursor + 1) % 2
			case "enter":
				if m.initialConfirmCursor == 0 {
					m.step = stepType
				} else {
					m.quitting = true
					return m, tea.Quit
				}
			}
		case stepType:
			switch msg.String() {
			case "up", "k":
				if m.commitTypeCursor > 0 {
					m.commitTypeCursor--
				}
			case "down", "j":
				if m.commitTypeCursor < len(m.commitTypeChoices)-1 {
					m.commitTypeCursor++
				}
			case "enter":
				m.step++
			}
		case stepScope:
			switch msg.String() {
			case "up", "k":
				if m.scopeCursor > 0 {
					m.scopeCursor--
				}
			case "down", "j":
				if m.scopeCursor < len(m.scopeChoices)-1 {
					m.scopeCursor++
				}
			case " ":
				if _, ok := m.selectedScopes[m.scopeCursor]; ok {
					delete(m.selectedScopes, m.scopeCursor)
				} else {
					m.selectedScopes[m.scopeCursor] = struct{}{}
				}
			case "enter":
				m.step++
				cmds = append(cmds, m.subjectInput.Focus())
			}
		case stepSubject:
			if msg.Type == tea.KeyEnter {
				if m.subjectInput.Value() == "" {
					m.errorMessage = "The commit title cannot be empty."
				} else {
					m.step++
					cmds = append(cmds, m.bodyInput.Focus())
				}
			} else {
				m.subjectInput, cmd = m.subjectInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		case stepBody:
			if msg.Type == tea.KeyCtrlD {
				m.step++
				m.bodyInput.Blur()
			} else {
				m.bodyInput, cmd = m.bodyInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		case stepBreakingChangePrompt:
			switch msg.String() {
			case "up", "k", "down", "j":
				m.breakingChangePromptCursor = (m.breakingChangePromptCursor + 1) % 2
			case "enter":
				if m.breakingChangePromptCursor == 1 {
					m.step++
					cmds = append(cmds, m.breakingChangeInput.Focus())
				} else {
					m.step = stepConfirm
				}
			}
		case stepBreakingChangeDescription:
			if msg.Type == tea.KeyCtrlD {
				if m.breakingChangeInput.Value() == "" {
					m.errorMessage = "The breaking change description cannot be empty."
				} else {
					m.step = stepConfirm
					m.breakingChangeInput.Blur()
				}
			} else {
				m.breakingChangeInput, cmd = m.breakingChangeInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		case stepConfirm:
			switch msg.String() {
			case "up", "k", "down", "j":
				m.confirmCursor = (m.confirmCursor + 1) % 2
			case "enter":
				if m.confirmCursor == 0 {
					m.quitting = true
					m.finalOutput = m.prepareHookOutput()
					return m, tea.Quit
				} else {
					return newModel(m.initialMessage, m.config), nil
				}
			}
		}
	}
	return m, tea.Batch(cmds...)
}
