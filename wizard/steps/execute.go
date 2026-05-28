package steps

import (
	"time"

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

// TaskProgressMsg is sent to update a task's state.
type TaskProgressMsg struct {
	Index int
	State TaskState
	Err   error
}

// SpinnerTickMsg advances the spinner animation frame.
type SpinnerTickMsg struct{}

// SpinnerTick returns a command that fires SpinnerTickMsg after 80ms.
func SpinnerTick() tea.Cmd {
	return tea.Tick(80*time.Millisecond, func(time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var (
	doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	runningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	failedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	pendingStyle = lipgloss.NewStyle().Faint(true)
)

type ExecuteStep struct {
	tasks        []TaskStatus
	spinnerFrame int
	done         bool
}

func NewExecuteStep(tasks []TaskStatus) Step {
	return &ExecuteStep{tasks: tasks}
}

func (s *ExecuteStep) Update(msg tea.Msg) (Step, tea.Cmd) {
	cp := *s
	cp.tasks = make([]TaskStatus, len(s.tasks))
	copy(cp.tasks, s.tasks)

	switch msg := msg.(type) {
	case TaskProgressMsg:
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
		// Keep spinner going while any task is still running
		if !cp.done {
			return &cp, SpinnerTick()
		}
	case SpinnerTickMsg:
		cp.spinnerFrame++
		// Continue ticking while any task is running
		for _, t := range cp.tasks {
			if t.State == TaskRunning {
				return &cp, SpinnerTick()
			}
		}
	}
	return &cp, nil
}

func (s *ExecuteStep) View() string {
	frame := spinnerFrames[s.spinnerFrame%len(spinnerFrames)]
	out := titleStyle.Render("Setting up your project") + "\n\n"
	for _, t := range s.tasks {
		var icon string
		switch t.State {
		case TaskDone:
			icon = doneStyle.Render("✓")
		case TaskRunning:
			icon = runningStyle.Render(frame)
		case TaskFailed:
			icon = failedStyle.Render("✗")
		default:
			icon = pendingStyle.Render("○")
		}
		out += icon + "  " + t.Label + "\n"
	}
	if s.done {
		out += "\n" + doneStyle.Bold(true).Render("✓ All done!")
	}
	return out
}

func (s *ExecuteStep) IsDone() bool  { return s.done }
func (s *ExecuteStep) Value() any    { return s.tasks }
func (s *ExecuteStep) Label() string { return "execute" }
