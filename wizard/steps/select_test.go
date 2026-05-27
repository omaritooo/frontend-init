package steps_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/omaritooo/frontend-init/wizard/steps"
	"github.com/stretchr/testify/assert"
)

func TestSelectStep_NavigatesDown(t *testing.T) {
	s := steps.NewSelectStep("Pick one", []string{"a", "b", "c"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
	assert.Equal(t, "b", result.Value())
}

func TestSelectStep_NavigatesUp(t *testing.T) {
	s := steps.NewSelectStep("Pick one", []string{"a", "b", "c"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
	result, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, "a", result.Value())
}

func TestSelectStep_DoesNotGoAboveFirst(t *testing.T) {
	s := steps.NewSelectStep("Pick one", []string{"a", "b"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyUp})
	assert.Equal(t, "a", result.Value())
}

func TestSelectStep_EnterCompletes(t *testing.T) {
	s := steps.NewSelectStep("Pick one", []string{"a", "b"})
	result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
	assert.True(t, result.IsDone())
	assert.Equal(t, "a", result.Value())
}

func TestSelectStep_ViewContainsTitle(t *testing.T) {
	s := steps.NewSelectStep("Choose framework", []string{"react", "vue"})
	assert.Contains(t, s.View(), "Choose framework")
}

func TestSelectStep_LabelReturnsTitle(t *testing.T) {
	s := steps.NewSelectStep("framework", []string{"react", "vue"})
	assert.Equal(t, "framework", s.Label())
}
