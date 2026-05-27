package wizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/executor"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

// stepLabelFramework is the label of the framework select step.
// Using a constant prevents a silent breakage if the label is ever renamed.
const stepLabelFramework = "framework"

// Model is the root Bubbletea model for the wizard.
type Model struct {
	stepList []steps.Step
	cursor   int
	cfg      *config.ProjectConfig
}

// New creates a new wizard Model with the full initial step list.
func New(cfg *config.ProjectConfig) Model {
	return Model{
		stepList: BuildInitialSteps(cfg),
		cfg:      cfg,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		}
	}

	current := m.stepList[m.cursor]
	newStep, cmd := current.Update(msg)

	// copy the step list to avoid mutating the backing array
	newList := make([]steps.Step, len(m.stepList))
	copy(newList, m.stepList)
	newList[m.cursor] = newStep
	m.stepList = newList

	if newStep.IsDone() {
		m.applyStepValue(newStep)
		if m.cursor < len(m.stepList)-1 {
			m.cursor++
			if newStep.Label() == stepLabelFramework {
				m.stepList = rebuildAfterFramework(m.stepList, m.cursor, m.cfg)
			}
		} else {
			return m, tea.Quit
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return m.stepList[m.cursor].View()
}

// Cursor returns the current step index (for testing).
func (m Model) Cursor() int { return m.cursor }

// Config returns the underlying ProjectConfig (for the executor).
func (m Model) Config() *config.ProjectConfig { return m.cfg }

// applyStepValue writes the completed step's result into ProjectConfig.
// It intentionally uses a value receiver: cfg is already a pointer, so mutations
// through m.cfg reach the underlying struct without needing a pointer receiver here.
func (m Model) applyStepValue(s steps.Step) {
	val := s.Value()
	switch s.Label() {
	case "mode":
		m.cfg.Mode = val.(string)
	case "preset":
		name := val.(string)
		if name != "Custom" {
			if p := config.GetPreset(name); p != nil {
				p.Apply(m.cfg)
			}
		}
		m.cfg.Preset = name
	case "package manager":
		m.cfg.PackageManager = val.(string)
	case stepLabelFramework:
		m.cfg.Framework = val.(string)
	case "variant":
		m.cfg.Variant = val.(string)
	case "typescript":
		m.cfg.TypeScript = val.(string) == "yes"
	case "linting":
		m.cfg.Linting = val.(string)
	case "ui library":
		m.cfg.UILibrary = val.(string)
	case "testing":
		m.cfg.Testing = val.([]string)
	case "tooling":
		m.cfg.Tooling = val.([]string)
	}
}

// rebuildAfterFramework replaces the step list from `from` onward with
// framework-appropriate steps (variant, typescript, linting, ui library,
// testing, tooling, confirm). Called after the framework step is completed.
func rebuildAfterFramework(current []steps.Step, from int, cfg *config.ProjectConfig) []steps.Step {
	head := make([]steps.Step, from)
	copy(head, current)

	tail := variantSteps(cfg)
	tail = append(tail,
		steps.NewSelectStep("typescript", []string{"yes", "no"}),
		steps.NewSelectStep("linting", []string{"eslint-prettier", "biome", "oxlint", "none"}),
		steps.NewSelectStep("ui library", UILibraryOptions(cfg)),
		steps.NewMultiSelectStep("testing", TestingOptions(cfg)),
		steps.NewMultiSelectStep("tooling", ToolingOptions(cfg)),
		steps.NewConfirmStep(cfg),
	)
	return append(head, tail...)
}

// ExecuteModel is the Bubbletea model for the execution progress screen.
type ExecuteModel struct {
	step  steps.Step
	tasks []executor.Task
	index int
}

// NewExecuteModel creates an ExecuteModel from an ExecuteStep and the executor task list.
func NewExecuteModel(step steps.Step, exTasks []executor.Task) ExecuteModel {
	return ExecuteModel{step: step, tasks: exTasks}
}

func (e ExecuteModel) Init() tea.Cmd {
	if len(e.tasks) == 0 {
		return tea.Quit
	}
	return runNextTask(e.tasks, 0)
}

func (e ExecuteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case steps.TaskProgressMsg:
		newStep, _ := e.step.Update(msg)
		e.step = newStep
		if msg.State == steps.TaskDone || msg.State == steps.TaskFailed {
			if e.index+1 < len(e.tasks) {
				e.index++
				return e, runNextTask(e.tasks, e.index)
			}
		}
		if e.step.IsDone() {
			return e, tea.Quit
		}
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return e, tea.Quit
		}
	}
	return e, nil
}

func (e ExecuteModel) View() string { return e.step.View() }

func runNextTask(tasks []executor.Task, idx int) tea.Cmd {
	return func() tea.Msg {
		err := tasks[idx].Run()
		state := steps.TaskDone
		if err != nil {
			state = steps.TaskFailed
		}
		return steps.TaskProgressMsg{Index: idx, State: state, Err: err}
	}
}
