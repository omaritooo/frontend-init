package steps

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskState int

const (
	TaskPending TaskState = iota
	TaskRunning
	TaskDone
	TaskFailed
)

type TaskStatus struct {
	Label string
	State TaskState
	Err   error
}

// TaskProgressMsg is sent by the executor to update a task's state.
type TaskProgressMsg struct {
	Index int
	State TaskState
	Err   error
}

var (
	doneIcon    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
	runningIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("⠋")
	pendingIcon = lipgloss.NewStyle().Faint(true).Render("○")
	failedIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗")
)

type ExecuteStep struct {
	tasks []TaskStatus
	done  bool
}

func NewExecuteStep(tasks []TaskStatus) Step {
	return &ExecuteStep{tasks: tasks}
}

func (s *ExecuteStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	cp := *s
	cp.tasks = make([]TaskStatus, len(s.tasks))
	copy(cp.tasks, s.tasks)

	if msg, ok := msg.(TaskProgressMsg); ok {
		if msg.Index >= 0 && msg.Index < len(cp.tasks) {
			cp.tasks[msg.Index].State = msg.State
			cp.tasks[msg.Index].Err = msg.Err
		}
		allDone := true
		for _, t := range cp.tasks {
			if t.State != TaskDone && t.State != TaskFailed {
				allDone = false
				break
			}
		}
		cp.done = allDone && len(cp.tasks) > 0
	}
	return &cp, nil
}

func (s *ExecuteStep) View() string {
	out := titleStyle.Render("Setting up your project") + "\n\n"
	for _, t := range s.tasks {
		icon := pendingIcon
		switch t.State {
		case TaskDone:
			icon = doneIcon
		case TaskRunning:
			icon = runningIcon
		case TaskFailed:
			icon = failedIcon
		}
		out += icon + "  " + t.Label + "\n"
	}
	if s.done {
		out += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render("✓ All done!")
	}
	return out
}

func (s *ExecuteStep) IsDone() bool  { return s.done }
func (s *ExecuteStep) Value() any    { return s.tasks }
func (s *ExecuteStep) Label() string { return "execute" }
