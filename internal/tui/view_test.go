package tui

import (
	"strings"
	"testing"

	"github.com/luizvilasboas/commit-hooks/internal/config"
)

func TestPrepareHookOutput(t *testing.T) {
	sampleConfig := config.Config{
		Types:  []string{"feat", "fix"},
		Scopes: []string{"api", "ui"},
	}

	t.Run("with type and subject", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.commitTypeCursor = 0 // feat
		m.subjectInput.SetValue("test subject")
		expected := "feat: test subject"
		if output := m.prepareHookOutput(); output != expected {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("with type, scope, and subject", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.commitTypeCursor = 1                     // fix
		m.selectedScopes = map[int]struct{}{0: {}} // api
		m.subjectInput.SetValue("bug in endpoint")
		expected := "fix(api): bug in endpoint"
		if output := m.prepareHookOutput(); output != expected {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("with multiple scopes", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.commitTypeCursor = 0                            // feat
		m.selectedScopes = map[int]struct{}{0: {}, 1: {}} // api, ui
		m.subjectInput.SetValue("new feature")
		expected := "feat(api, ui): new feature"
		if output := m.prepareHookOutput(); output != expected {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("with body", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.commitTypeCursor = 0 // feat
		m.subjectInput.SetValue("test subject")
		m.bodyInput.SetValue("This is a detailed body.")
		expected := "feat: test subject\n\nThis is a detailed body."
		if output := m.prepareHookOutput(); output != expected {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("with breaking change", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.commitTypeCursor = 0 // feat
		m.subjectInput.SetValue("test subject")
		m.breakingChangeInput.SetValue("The API now requires authentication.")
		expected := "feat: test subject\n\nBREAKING CHANGE: The API now requires authentication."
		if output := m.prepareHookOutput(); output != expected {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})

	t.Run("with all parts", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.commitTypeCursor = 0                     // feat
		m.selectedScopes = map[int]struct{}{0: {}} // api
		m.subjectInput.SetValue("refactor login")
		m.bodyInput.SetValue("The login function was refactored.")
		m.breakingChangeInput.SetValue("The old password hash is no longer supported.")
		expected := "feat(api): refactor login\n\nThe login function was refactored.\n\nBREAKING CHANGE: The old password hash is no longer supported."
		if output := m.prepareHookOutput(); output != expected {
			t.Errorf("expected %q, got %q", expected, output)
		}
	})
}

func TestView(t *testing.T) {
	sampleConfig := config.Config{
		Types:  []string{"feat", "fix"},
		Scopes: []string{"api", "ui"},
	}

	t.Run("should render initial confirm view", func(t *testing.T) {
		m := newModel("initial message", sampleConfig)
		view := m.View()

		if !strings.Contains(view, "You provided the following commit message:") {
			t.Error("view does not contain initial confirm prompt")
		}
		if !strings.Contains(view, "initial message") {
			t.Error("view does not contain the initial message")
		}
		if !strings.Contains(view, "> Yes, enhance this message") {
			t.Error("view does not show the 'Yes' option with a cursor")
		}
	})

	t.Run("should render type select view", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepType
		m.commitTypeCursor = 1 // fix
		view := m.View()

		if !strings.Contains(view, "What is the type of the commit?") {
			t.Error("view does not contain type select prompt")
		}
		if !strings.Contains(view, "> fix") {
			t.Error("view does not show the 'fix' option with a cursor")
		}
	})

	t.Run("should render scope checkbox view", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepScope
		m.scopeCursor = 0                // api
		m.selectedScopes[1] = struct{}{} // ui is selected
		view := m.View()

		if !strings.Contains(view, "What is the scope of this change?") {
			t.Error("view does not contain scope checkbox prompt")
		}
		if !strings.Contains(view, "> [ ] api") {
			t.Error("view does not show cursor on unselected 'api' option")
		}
		if !strings.Contains(view, "[x] ui") {
			t.Error("view does not show 'ui' as selected")
		}
	})

	t.Run("should render subject input view", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepSubject
		view := m.View()
		if !strings.Contains(view, "Write a short, imperative tense title:") {
			t.Error("view does not contain subject input prompt")
		}
	})

	t.Run("should render completed steps", func(t *testing.T) {
		m := newModel("initial", sampleConfig)
		m.step = stepSubject             // At subject step
		m.commitTypeCursor = 0           // feat
		m.selectedScopes[1] = struct{}{} // ui
		view := m.View()

		if !strings.Contains(view, "Original Message:\n> initial") {
			t.Error("view does not show the original message")
		}
		if !strings.Contains(view, "What is the type of the commit?\n> feat") {
			t.Error("view does not show the completed type selection")
		}
		if !strings.Contains(view, "What is the scope of this change? (optional)\n> ui") {
			t.Error("view does not show the completed scope selection")
		}
	})

	t.Run("should render error message", func(t *testing.T) {
		m := newModel("", sampleConfig)
		m.step = stepSubject
		m.errorMessage = "This is an error"
		view := m.View()

		if !strings.Contains(view, "This is an error") {
			t.Error("view does not display the error message")
		}
	})

}
