package wizard

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/executor"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

const stepLabelFramework = "framework"
const stepLabelMode = "mode"

type Model struct {
	stepList []steps.Step
	cursor   int
	cfg      *config.ProjectConfig
}

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

	newList := make([]steps.Step, len(m.stepList))
	copy(newList, m.stepList)
	newList[m.cursor] = newStep
	m.stepList = newList

	if newStep.IsDone() {
		m.applyStepValue(newStep)
		if m.cursor < len(m.stepList)-1 {
			m.cursor++
			switch newStep.Label() {
			case stepLabelMode:
				if m.cfg.Mode == "new" {
					m.stepList = insertStepAt(m.stepList, m.cursor, steps.NewInputStep("project name", "my-app"))
				}
			case stepLabelFramework:
				m.stepList = rebuildAfterFramework(m.stepList, m.cursor, m.cfg)
			case "ui library":
				if m.cfg.UILibrary == "shadcn" || m.cfg.UILibrary == "shadcn-svelte" {
					m.stepList = insertStepAt(m.stepList, m.cursor, steps.NewSelectStep("shadcn theme", shadcnThemes()))
				}
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

func (m Model) Cursor() int { return m.cursor }

func (m Model) Config() *config.ProjectConfig { return m.cfg }

func (m Model) applyStepValue(s steps.Step) {
	val := s.Value()
	switch s.Label() {
	case stepLabelMode:
		m.cfg.Mode = val.(string)
	case "project name":
		m.cfg.ProjectName = val.(string)
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
	case "shadcn theme":
		m.cfg.ShadcnTheme = val.(string)
	case "testing":
		m.cfg.Testing = val.([]string)
	case "tooling":
		m.cfg.Tooling = val.([]string)
	}
}

func insertStepAt(list []steps.Step, idx int, s steps.Step) []steps.Step {
	out := make([]steps.Step, len(list)+1)
	copy(out, list[:idx])
	out[idx] = s
	copy(out[idx+1:], list[idx:])
	return out
}

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

type ExecuteModel struct {
	step  steps.Step
	tasks []executor.Task
	index int
}

func NewExecuteModel(step steps.Step, exTasks []executor.Task) ExecuteModel {
	return ExecuteModel{step: step, tasks: exTasks}
}

func (e ExecuteModel) Init() tea.Cmd {
	if len(e.tasks) == 0 {
		return tea.Quit
	}
	return markRunning(0)
}

func (e ExecuteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case steps.TaskProgressMsg:
		newStep, stepCmd := e.step.Update(msg)
		e.step = newStep
		switch msg.State {
		case steps.TaskRunning:
			// Now actually execute the task in a background goroutine.
			return e, tea.Batch(stepCmd, executeTask(e.tasks, msg.Index))
		case steps.TaskDone, steps.TaskFailed:
			if e.index+1 < len(e.tasks) {
				e.index++
				next := e.index
				return e, tea.Batch(stepCmd, markRunning(next))
			}
		}
		if e.step.IsDone() {
			return e, tea.Quit
		}
		return e, stepCmd
	case steps.SpinnerTickMsg:
		newStep, stepCmd := e.step.Update(msg)
		e.step = newStep
		return e, stepCmd
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return e, tea.Quit
		}
	}
	return e, nil
}

func (e ExecuteModel) View() string { return e.step.View() }

func markRunning(idx int) tea.Cmd {
	return func() tea.Msg {
		return steps.TaskProgressMsg{Index: idx, State: steps.TaskRunning}
	}
}

func executeTask(tasks []executor.Task, idx int) tea.Cmd {
	return func() tea.Msg {
		err := tasks[idx].Run()
		state := steps.TaskDone
		if err != nil {
			state = steps.TaskFailed
		}
		return steps.TaskProgressMsg{Index: idx, State: state, Err: err}
	}
}
