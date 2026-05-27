package steps_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

func TestMultiSelectStep_ToggleSelection(t *testing.T) {
	s := steps.NewMultiSelectStep("Pick tools", []string{"vitest", "playwright", "storybook"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeySpace})
	vals := result.Value().([]string)
	assert.Contains(t, vals, "vitest")
}

func TestMultiSelectStep_DeselectOnSecondToggle(t *testing.T) {
	s := steps.NewMultiSelectStep("Pick tools", []string{"vitest", "playwright"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeySpace})
	result, _ = result.Update(tea.KeyMsg{Type: tea.KeySpace})
	vals := result.Value().([]string)
	assert.NotContains(t, vals, "vitest")
}

func TestMultiSelectStep_EnterCompletes(t *testing.T) {
	s := steps.NewMultiSelectStep("Pick tools", []string{"vitest"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, result.IsDone())
}

func TestMultiSelectStep_CanSelectNone(t *testing.T) {
	s := steps.NewMultiSelectStep("Pick tools", []string{"vitest"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	vals := result.Value().([]string)
	assert.Empty(t, vals)
}

func TestMultiSelectStep_NavigatesDown(t *testing.T) {
	s := steps.NewMultiSelectStep("Pick tools", []string{"a", "b", "c"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
	// select item at cursor (now "b")
	result, _ = result.Update(tea.KeyMsg{Type: tea.KeySpace})
	vals := result.Value().([]string)
	assert.Contains(t, vals, "b")
	assert.NotContains(t, vals, "a")
}

func TestMultiSelectStep_LabelReturnsTitle(t *testing.T) {
	s := steps.NewMultiSelectStep("tooling", []string{"zustand"})
	assert.Equal(t, "tooling", s.Label())
}
