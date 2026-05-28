package steps

import (
	tea "github.com/charmbracelet/bubbletea"
)

type InputStep struct {
	title       string
	placeholder string
	value       string
	done        bool
}

func NewInputStep(title, placeholder string) Step {
	return &InputStep{title: title, placeholder: placeholder}
}

func (s *InputStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	cp := *s
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyEnter:
			if cp.value == "" {
				cp.value = cp.placeholder
			}
			cp.done = true
		case tea.KeyBackspace, tea.KeyDelete:
			if len(cp.value) > 0 {
				cp.value = cp.value[:len(cp.value)-1]
			}
		case tea.KeyRunes:
			cp.value += string(key.Runes)
		}
	}
	return &cp, nil
}

func (s *InputStep) View() string {
	out := titleStyle.Render(s.title) + "\n\n"
	display := s.value
	if display == "" {
		display = hintStyle.Render(s.placeholder)
	}
	out += "  " + display + cursorStyle.Render("▌") + "\n"
	out += "\n" + hintStyle.Render("type • backspace to delete • enter to confirm")
	return out
}

func (s *InputStep) IsDone() bool  { return s.done }
func (s *InputStep) Value() any    { return s.value }
func (s *InputStep) Label() string { return s.title }
