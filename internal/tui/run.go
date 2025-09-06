package tui

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luizvilasboas/commit-hooks/internal/config"
)

func Run(initialMsg string, cfg config.Config) (string, error) {
	p := tea.NewProgram(newModel(initialMsg, cfg), tea.WithOutput(os.Stderr))

	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(model); ok && m.finalOutput != "" {
		return m.finalOutput, nil
	}

	return "", nil
}
