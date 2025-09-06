package tui

import (
	"testing"

	"github.com/luizvilasboas/commit-hooks/internal/config"

	tea "github.com/charmbracelet/bubbletea"
)

func key(k string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
}

var keyEnter = tea.KeyMsg{Type: tea.KeyEnter}

func TestUpdate(t *testing.T) {
	sampleConfig := config.Config{
		Types:  []string{"feat", "fix"},
		Scopes: []string{"api", "ui"},
	}

	t.Run("should quit on Ctrl+C", func(t *testing.T) {
		m := newModel("", sampleConfig)
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}

		newModel, _ := m.Update(msg)

		if !newModel.(model).quitting {
			t.Error("expected quitting to be true on Ctrl+C")
		}
	})

	t.Run("stepInitialConfirm - should advance to stepType on Yes", func(t *testing.T) {
		m := newModel("initial", sampleConfig)
		m.initialConfirmCursor = 0

		newModel, _ := m.Update(keyEnter)

		if newModel.(model).step != stepType {
			t.Errorf("expected step to be %d, got %d", stepType, newModel.(model).step)
		}
	})

	t.Run("stepInitialConfirm - should quit on No", func(t *testing.T) {
		m := newModel("initial", sampleConfig)
		m.initialConfirmCursor = 1

		newModel, _ := m.Update(keyEnter)

		if !newModel.(model).quitting {
			t.Error("expected quitting to be true when selecting No")
		}
	})

	t.Run("stepType - should advance to stepScope on Enter", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepType

		newModel, _ := m.Update(keyEnter)

		if newModel.(model).step != stepScope {
			t.Errorf("expected step to be %d, got %d", stepScope, newModel.(model).step)
		}
	})

	t.Run("stepScope - should select/deselect scope on Space", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepScope
		m.scopeCursor = 1

		updatedModel, _ := m.Update(key(" "))
		m = updatedModel.(model)
		if _, ok := m.selectedScopes[1]; !ok {
			t.Error("expected scope at cursor 1 to be selected")
		}

		updatedModel, _ = m.Update(key(" "))
		m = updatedModel.(model)
		if _, ok := m.selectedScopes[1]; ok {
			t.Error("expected scope at cursor 1 to be deselected")
		}
	})

	t.Run("stepSubject - should show error on empty subject", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepSubject
		m.subjectInput.SetValue("")

		newModel, _ := m.Update(keyEnter)

		if newModel.(model).errorMessage == "" {
			t.Error("expected an error message for empty subject")
		}
		if newModel.(model).step != stepSubject {
			t.Error("expected to remain on stepSubject after error")
		}
	})

	t.Run("stepSubject - should advance on valid subject", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepSubject
		m.subjectInput.SetValue("valid subject")

		newModel, _ := m.Update(keyEnter)

		if newModel.(model).step != stepBody {
			t.Errorf("expected step to be %d, got %d", stepBody, newModel.(model).step)
		}
	})

	t.Run("stepConfirm - should reset on Start Over", func(t *testing.T) {
		m := newModel("initial", sampleConfig)
		m.step = stepConfirm
		m.confirmCursor = 1

		newModel, _ := m.Update(keyEnter)

		if newModel.(model).step != stepInitialConfirm {
			t.Errorf("expected step to be reset to %d, got %d", stepInitialConfirm, newModel.(model).step)
		}
	})

	t.Run("stepConfirm - should quit and set finalOutput on Confirm", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepConfirm
		m.commitTypeCursor = 0
		m.subjectInput.SetValue("test subject")
		m.confirmCursor = 0

		newModel, _ := m.Update(keyEnter)

		finalM := newModel.(model)
		if !finalM.quitting {
			t.Error("expected quitting to be true on confirm")
		}
		if finalM.finalOutput != "feat: test subject" {
			t.Errorf("expected finalOutput to be 'feat: test subject', got %q", finalM.finalOutput)
		}
	})
}
