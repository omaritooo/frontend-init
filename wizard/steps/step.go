package steps

import tea "github.com/charmbracelet/bubbletea"

// Step is the interface every wizard screen implements.
type Step interface {
	Update(tea.Msg) (Step, tea.Cmd)
	View() string
	IsDone() bool
	Value() any
	Label() string // lowercase title identifying this step
}
