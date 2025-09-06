package tui

import (
	"fmt"
	"strings"
)

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	if m.step == stepInitialConfirm {
		b.WriteString(m.viewInitialConfirm())
	} else {
		if m.initialMessage != "" {
			b.WriteString(fmt.Sprintf("%s\n%s\n\n", "Original Message:", "> "+m.initialMessage))
		}
		if m.step > stepType {
			b.WriteString(m.viewCompletedType())
		}
		if m.step > stepScope {
			b.WriteString(m.viewCompletedScope())
		}
	}

	switch m.step {
	case stepType:
		b.WriteString(m.viewTypeSelect())
	case stepScope:
		b.WriteString(m.viewScopeCheckbox())
	case stepSubject:
		b.WriteString(m.viewSubjectInput())
	case stepBody:
		b.WriteString(m.viewBodyInput())
	case stepBreakingChangePrompt:
		b.WriteString(m.viewBreakingChangePrompt())
	case stepBreakingChangeDescription:
		b.WriteString(m.viewBreakingChangeDescriptionInput())
	case stepConfirm:
		b.WriteString(m.viewConfirm())
	}

	if m.errorMessage != "" {
		b.WriteString("\n" + m.errorMessage)
	}

	return b.String()
}

func (m model) viewInitialConfirm() string {
	var b strings.Builder
	b.WriteString("You provided the following commit message:" + "\n")
	b.WriteString(m.initialMessage + "\n\n")
	b.WriteString("Do you want to use the helper to enhance it into a semantic commit?" + "\n\n")

	for i, choice := range m.initialConfirmChoices {
		cursor := " "
		if m.initialConfirmCursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}

func (m model) prepareHookOutput() string {
	var b strings.Builder
	commitType := m.commitTypeChoices[m.commitTypeCursor]
	subject := m.subjectInput.Value()
	body := m.bodyInput.Value()
	breakingChange := m.breakingChangeInput.Value()
	var selectedScopesStr []string
	for i, choice := range m.scopeChoices {
		if _, ok := m.selectedScopes[i]; ok {
			selectedScopesStr = append(selectedScopesStr, choice)
		}
	}
	scope := strings.Join(selectedScopesStr, ", ")
	if scope != "" {
		b.WriteString(fmt.Sprintf("%s(%s): %s", commitType, scope, subject))
	} else {
		b.WriteString(fmt.Sprintf("%s: %s", commitType, subject))
	}
	if body != "" {
		b.WriteString(fmt.Sprintf("\n\n%s", body))
	}
	if breakingChange != "" {
		b.WriteString(fmt.Sprintf("\n\nBREAKING CHANGE: %s", breakingChange))
	}
	return b.String()
}

func (m model) viewConfirm() string {
	var b strings.Builder
	b.WriteString("Commit Review:" + "\n")
	b.WriteString("-------------------\n")
	b.WriteString(m.prepareHookOutput())
	b.WriteString("\n-------------------\n\n")
	for i, choice := range m.confirmChoices {
		cursor := " "
		if m.confirmCursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}

func (m model) viewTypeSelect() string {
	var b strings.Builder
	b.WriteString("What is the type of the commit?" + "\n\n")
	for i, choice := range m.commitTypeChoices {
		cursor := " "
		if m.commitTypeCursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	b.WriteString("\n" + "Use arrow keys and Enter to select.")
	return b.String()
}

func (m model) viewScopeCheckbox() string {
	var b strings.Builder
	b.WriteString("What is the scope of this change? (optional)" + "\n\n")
	for i, choice := range m.scopeChoices {
		cursor := " "
		if m.scopeCursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.selectedScopes[i]; ok {
			checked = "x"
		}
		b.WriteString(fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice))
	}
	b.WriteString("\n" + "Use arrows to move, Space to select, Enter to continue.")
	return b.String()
}

func (m model) viewSubjectInput() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", "Write a short, imperative tense title:", m.subjectInput.View(), "Press Enter to continue.")
}

func (m model) viewBodyInput() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", "Add a longer description: (optional)", m.bodyInput.View(), "Press Ctrl+D to finish/skip this step.")
}

func (m model) viewBreakingChangePrompt() string {
	var b strings.Builder
	b.WriteString("Is it a breaking change?" + "\n\n")
	for i, choice := range m.breakingChangePrompt {
		cursor := " "
		if m.breakingChangePromptCursor == i {
			cursor = ">"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
	return b.String()
}

func (m model) viewBreakingChangeDescriptionInput() string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", "Describe the breaking change:", m.breakingChangeInput.View(), "Press Ctrl+D to finish.")
}

func (m model) viewCompletedType() string {
	return fmt.Sprintf("%s\n%s\n\n", "What is the type of the commit?", "> "+m.commitTypeChoices[m.commitTypeCursor])
}

func (m model) viewCompletedScope() string {
	var selected []string
	for i, choice := range m.scopeChoices {
		if _, ok := m.selectedScopes[i]; ok {
			selected = append(selected, choice)
		}
	}
	answer := "> (none)"
	if len(selected) > 0 {
		answer = "> " + strings.Join(selected, ", ")
	}
	return fmt.Sprintf("%s\n%s\n\n", "What is the scope of this change? (optional)", answer)
}
