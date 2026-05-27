package steps

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	titleStyle    = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	hintStyle     = lipgloss.NewStyle().Faint(true)
)

type SelectStep struct {
	title   string
	choices []string
	cursor  int
	done    bool
}

func NewSelectStep(title string, choices []string) Step {
	return &SelectStep{title: title, choices: choices}
}

func (s *SelectStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	cp := *s
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if cp.cursor > 0 {
				cp.cursor--
			}
		case tea.KeyDown:
			if cp.cursor < len(cp.choices)-1 {
				cp.cursor++
			}
		case tea.KeyEnter:
			cp.done = true
		}
	}
	return &cp, nil
}

func (s *SelectStep) View() string {
	out := titleStyle.Render(s.title) + "\n\n"
	for i, c := range s.choices {
		cursor := "  "
		line := c
		if i == s.cursor {
			cursor = cursorStyle.Render("▶ ")
			line = selectedStyle.Render(c)
		}
		out += fmt.Sprintf("%s%s\n", cursor, line)
	}
	out += "\n" + hintStyle.Render("↑/↓ navigate • enter select")
	return out
}

func (s *SelectStep) IsDone() bool  { return s.done }
func (s *SelectStep) Value() any    { return s.choices[s.cursor] }
func (s *SelectStep) Label() string { return s.title }
