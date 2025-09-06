package tui

import (
	"reflect"
	"testing"

	"github.com/luizvilasboas/commit-hooks/internal/config"
)

func TestNewModel(t *testing.T) {
	sampleConfig := config.Config{
		Types:  []string{"feat", "fix"},
		Scopes: []string{"api", "ui"},
	}

	t.Run("should initialize with initial message", func(t *testing.T) {
		initialMsg := "initial commit"
		m := newModel(initialMsg, sampleConfig)

		if m.step != stepInitialConfirm {
			t.Errorf("expected step to be %d, got %d", stepInitialConfirm, m.step)
		}

		if m.initialMessage != initialMsg {
			t.Errorf("expected initialMessage to be %q, got %q", initialMsg, m.initialMessage)
		}

		if m.subjectInput.Value() != initialMsg {
			t.Errorf("expected subjectInput value to be %q, got %q", initialMsg, m.subjectInput.Value())
		}

		if !reflect.DeepEqual(m.commitTypeChoices, sampleConfig.Types) {
			t.Errorf("expected commitTypeChoices to be %+v, got %+v", sampleConfig.Types, m.commitTypeChoices)
		}

		if !reflect.DeepEqual(m.scopeChoices, sampleConfig.Scopes) {
			t.Errorf("expected scopeChoices to be %+v, got %+v", sampleConfig.Scopes, m.scopeChoices)
		}
	})

	t.Run("should initialize without initial message", func(t *testing.T) {
		initialMsg := ""
		m := newModel(initialMsg, sampleConfig)

		if m.step != stepType {
			t.Errorf("expected step to be %d, got %d", stepType, m.step)
		}

		if m.initialMessage != initialMsg {
			t.Errorf("expected initialMessage to be empty, got %q", m.initialMessage)
		}

		if m.subjectInput.Value() != initialMsg {
			t.Errorf("expected subjectInput value to be empty, got %q", m.subjectInput.Value())
		}

		if !reflect.DeepEqual(m.commitTypeChoices, sampleConfig.Types) {
			t.Errorf("expected commitTypeChoices to be %+v, got %+v", sampleConfig.Types, m.commitTypeChoices)
		}
	})
}
