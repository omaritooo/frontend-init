package steps

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/omaritooo/frontend-init/config"
)

var (
	labelStyle = lipgloss.NewStyle().Faint(true).Width(18)
	valueStyle = lipgloss.NewStyle().Bold(true)
)

type ConfirmStep struct {
	cfg  *config.ProjectConfig
	done bool
}

func NewConfirmStep(cfg *config.ProjectConfig) Step {
	return &ConfirmStep{cfg: cfg}
}

func (s *ConfirmStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	cp := *s
	if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
		cp.done = true
	}
	return &cp, nil
}

func (s *ConfirmStep) View() string {
	row := func(label, val string) string {
		return labelStyle.Render(label+":") + " " + valueStyle.Render(val) + "\n"
	}
	out := titleStyle.Render("Review your setup") + "\n\n"
	out += row("Mode", orStr(s.cfg.Mode, "—"))
	out += row("Package manager", s.cfg.PackageManager)
	out += row("Framework", fmt.Sprintf("%s (%s)", s.cfg.Framework, s.cfg.Variant))
	out += row("TypeScript", fmt.Sprintf("%v", s.cfg.TypeScript))
	out += row("Linting", orStr(s.cfg.Linting, "none"))
	out += row("UI library", orStr(s.cfg.UILibrary, "none"))
	out += row("Testing", orStr(strings.Join(s.cfg.Testing, ", "), "none"))
	out += row("Tooling", orStr(strings.Join(s.cfg.Tooling, ", "), "none"))
	out += "\n" + hintStyle.Render("enter to confirm • esc to go back")
	return out
}

func (s *ConfirmStep) IsDone() bool  { return s.done }
func (s *ConfirmStep) Value() any    { return true }
func (s *ConfirmStep) Label() string { return "confirm" }

func orStr(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
