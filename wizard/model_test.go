package wizard_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/wizard"
	"github.com/stretchr/testify/assert"
)

func TestModel_AdvancesOnStepComplete(t *testing.T) {
	cfg := config.New()
	m := wizard.New(cfg)
	assert.Equal(t, 0, m.Cursor())
	// press enter on first step
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	wm := newM.(wizard.Model)
	assert.Equal(t, 1, wm.Cursor())
}

func TestModel_BacktrackOnEsc(t *testing.T) {
	cfg := config.New()
	m := wizard.New(cfg)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // advance to step 1
	wm := newM.(wizard.Model)
	newM, _ = wm.Update(tea.KeyMsg{Type: tea.KeyEsc}) // go back
	wm = newM.(wizard.Model)
	assert.Equal(t, 0, wm.Cursor())
}

func TestModel_DoesNotGoBeforeFirst(t *testing.T) {
	cfg := config.New()
	m := wizard.New(cfg)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	wm := newM.(wizard.Model)
	assert.Equal(t, 0, wm.Cursor())
}

// TestModel_RebuildAfterFramework verifies that after completing the framework
// step the step list is rebuilt: cursor lands on a "variant" step and the view
// contains the react-specific choices (cursor starts at 0 → "react" is chosen).
func TestModel_RebuildAfterFramework(t *testing.T) {
	cfg := config.New()
	m := wizard.New(cfg)

	// Mode "new" causes a project name step to be inserted after mode.
	// Steps: mode → project name → preset → package manager → framework → (rebuild) → variant
	// 5 enters are needed to reach variant.
	var tm tea.Model = m
	for i := 0; i < 5; i++ {
		var cmd tea.Cmd
		tm, cmd = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
		_ = cmd
	}

	wm := tm.(wizard.Model)

	// After 5 enters the cursor should be at index 5 (the new "variant" step).
	assert.Equal(t, 5, wm.Cursor())

	// The current step should be the variant step injected by rebuildAfterFramework.
	view := wm.View()
	assert.Contains(t, view, "vite", "expected react variant choices to be visible")

	// Config.Framework must have been set to "react" (index 0 of the framework choices).
	assert.Equal(t, "react", wm.Config().Framework)
}

// TestModel_AppliesMultiSelectValue verifies that completing a multiselect step
// writes the result (an empty []string when nothing is toggled) into the config.
func TestModel_AppliesMultiSelectValue(t *testing.T) {
	cfg := config.New()
	m := wizard.New(cfg)

	// Navigate through 9 select steps:
	// mode → project name → preset → package manager → framework →
	// variant (after rebuild) → typescript → linting → ui library
	var tm tea.Model = m
	for i := 0; i < 9; i++ {
		var cmd tea.Cmd
		tm, cmd = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
		_ = cmd
	}

	// Now on step 9: the "testing" multiselect.
	// Press Enter without toggling anything → Value() returns []string{}.
	var cmd tea.Cmd
	tm, cmd = tm.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = cmd

	wm := tm.(wizard.Model)

	// Testing field must be populated (empty slice, not nil).
	assert.Equal(t, []string{}, wm.Config().Testing)
}
