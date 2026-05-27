package steps

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MultiSelectStep struct {
	title    string
	choices  []string
	selected map[int]bool
	cursor   int
	done     bool
}

func NewMultiSelectStep(title string, choices []string) Step {
	return &MultiSelectStep{
		title:    title,
		choices:  choices,
		selected: make(map[int]bool),
	}
}

func (s *MultiSelectStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	cp := *s
	cp.selected = make(map[int]bool, len(s.selected))
	for k, v := range s.selected {
		cp.selected[k] = v
	}
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
		case tea.KeySpace:
			cp.selected[cp.cursor] = !cp.selected[cp.cursor]
		case tea.KeyEnter:
			cp.done = true
		}
	}
	return &cp, nil
}

func (s *MultiSelectStep) View() string {
	out := titleStyle.Render(s.title) + "\n\n"
	for i, c := range s.choices {
		cur := "  "
		if i == s.cursor {
			cur = cursorStyle.Render("▶ ")
		}
		checkbox := "○"
		if s.selected[i] {
			checkbox = selectedStyle.Render("●")
		}
		out += fmt.Sprintf("%s%s %s\n", cur, checkbox, c)
	}
	out += "\n" + lipgloss.NewStyle().Faint(true).Render("↑/↓ navigate • space toggle • enter confirm")
	return out
}

func (s *MultiSelectStep) IsDone() bool { return s.done }

func (s *MultiSelectStep) Value() any {
	var result []string
	for i, c := range s.choices {
		if s.selected[i] {
			result = append(result, c)
		}
	}
	if result == nil {
		return []string{}
	}
	return result
}

func (s *MultiSelectStep) Label() string { return s.title }
