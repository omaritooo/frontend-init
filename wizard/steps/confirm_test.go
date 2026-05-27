package steps_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

func TestConfirmStep_ViewShowsAllChoices(t *testing.T) {
	cfg := config.New()
	cfg.Framework = "react"
	cfg.Variant   = "vite"
	cfg.Linting   = "eslint-prettier"
	cfg.UILibrary = "shadcn"
	cfg.Testing   = []string{"vitest", "playwright"}
	cfg.Tooling   = []string{"tanstack-query"}
	s := steps.NewConfirmStep(cfg)
	view := s.View()
	assert.Contains(t, view, "react")
	assert.Contains(t, view, "shadcn")
	assert.Contains(t, view, "vitest")
	assert.Contains(t, view, "tanstack-query")
}

func TestConfirmStep_EnterCompletes(t *testing.T) {
	s := steps.NewConfirmStep(config.New())
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, result.IsDone())
	assert.Equal(t, true, result.Value())
}

func TestConfirmStep_LabelIsConfirm(t *testing.T) {
	s := steps.NewConfirmStep(config.New())
	assert.Equal(t, "confirm", s.Label())
}
