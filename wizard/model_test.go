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
