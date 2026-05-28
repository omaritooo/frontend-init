package executor

import (
	"os/exec"
	"strings"

	"github.com/omaritooo/frontend-init/config"
	wsteps "github.com/omaritooo/frontend-init/wizard/steps"
)

// Task is a labelled unit of work with an execution function.
type Task struct {
	Label string
	Run   func() error
}

// Executor builds and runs the setup task list.
type Executor struct {
	cfg         *config.ProjectConfig
	runner      CommandRunner
	projectDir  string
	projectName string
}

// New creates an Executor for the given config, runner, and project directory.
func New(cfg *config.ProjectConfig, runner CommandRunner, projectDir string) *Executor {
	return &Executor{cfg: cfg, runner: runner, projectDir: projectDir}
}

// SetProjectName sets the project name used by scaffold commands.
func (e *Executor) SetProjectName(name string) { e.projectName = name }

// Tasks returns the ordered list of setup tasks derived from config.
func (e *Executor) Tasks() []Task {
	var tasks []Task

	if e.cfg.IsNewProject() {
		cmd := ScaffoldCommand(e.cfg, e.projectName)
		if cmd != nil {
			c := cmd
			tasks = append(tasks, Task{
				Label: "Scaffold project",
				Run:   func() error { return e.runner.Run(".", c[0], c[1:]...) },
			})
		}
	}

	tools := e.selectedTools()

	var pkgs, devPkgs []string
	for _, t := range tools {
		pkgs = append(pkgs, t.Packages...)
		devPkgs = append(devPkgs, t.DevPackages...)
	}
	if len(pkgs) > 0 {
		p := pkgs
		tasks = append(tasks, Task{
			Label: "Install dependencies",
			Run: func() error {
				cmd := InstallCmd(e.cfg.PackageManager, false, p)
				return e.runner.Run(e.projectDir, cmd[0], cmd[1:]...)
			},
		})
	}
	if len(devPkgs) > 0 {
		d := devPkgs
		tasks = append(tasks, Task{
			Label: "Install dev dependencies",
			Run: func() error {
				cmd := InstallCmd(e.cfg.PackageManager, true, d)
				return e.runner.Run(e.projectDir, cmd[0], cmd[1:]...)
			},
		})
	}

	for _, t := range tools {
		if len(t.ConfigFiles) > 0 {
			tool := t
			tasks = append(tasks, Task{
				Label: "Configure " + tool.Name,
				Run:   func() error { return WriteConfigFiles(e.projectDir, tool.ConfigFiles) },
			})
		}
	}

	for _, t := range tools {
		if len(t.FilePatches) > 0 {
			tool := t
			tasks = append(tasks, Task{
				Label: "Patch files for " + tool.Name,
				Run:   func() error { return ApplyPatches(e.projectDir, tool.FilePatches) },
			})
		}
	}

	allScripts := make(map[string]string)
	for _, t := range tools {
		for k, v := range t.Scripts {
			allScripts[k] = v
		}
	}
	if len(allScripts) > 0 {
		s := allScripts
		tasks = append(tasks, Task{
			Label: "Update package.json scripts",
			Run:   func() error { return MergeScripts(e.projectDir, s) },
		})
	}

	for _, t := range tools {
		for _, postCmd := range t.PostInstallCmds {
			pc := postCmd
			tool := t
			args := strings.Fields(pc)
			tasks = append(tasks, Task{
				Label: tool.Name + ": " + pc,
				Run: func() error {
					return e.runner.Run(e.projectDir, args[0], args[1:]...)
				},
			})
		}
	}

	return tasks
}

// selectedTools returns ToolSetup instances for all selected tools in config.
func (e *Executor) selectedTools() []ToolSetup {
	var tools []ToolSetup
	add := func(key string) {
		if s := GetToolSetup(key, e.cfg); s != nil {
			tools = append(tools, *s)
		}
	}
	if e.cfg.Linting != "none" && e.cfg.Linting != "" {
		add(e.cfg.Linting)
	}
	if e.cfg.UILibrary != "none" && e.cfg.UILibrary != "" {
		if e.cfg.UILibrary == "tailwind-only" {
			add("tailwind")
		} else {
			add("tailwind")
			add(e.cfg.UILibrary)
		}
	}
	for _, t := range e.cfg.Testing {
		add(t)
	}
	for _, t := range e.cfg.Tooling {
		add(t)
	}
	return tools
}

// RealRunner implements CommandRunner using os/exec.
type RealRunner struct{}

func (r *RealRunner) Run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}

// ToWizardTasks converts executor Tasks to wizard ExecuteStep TaskStatus slice.
func ToWizardTasks(tasks []Task) []wsteps.TaskStatus {
	result := make([]wsteps.TaskStatus, len(tasks))
	for i, t := range tasks {
		result[i] = wsteps.TaskStatus{Label: t.Label, State: wsteps.TaskPending}
	}
	return result
}
