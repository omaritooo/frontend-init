package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/executor"
	"github.com/omaritooo/frontend-init/wizard"
	wsteps "github.com/omaritooo/frontend-init/wizard/steps"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Start the interactive setup wizard",
	RunE:  runInit,
}

func runInit(_ *cobra.Command, _ []string) error {
	cfg := config.New()

	if wd, err := os.Getwd(); err == nil {
		cfg.PackageManager = executor.DetectPackageManager(wd)
	}

	m := wizard.New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	wm, ok := finalModel.(wizard.Model)
	if !ok {
		return fmt.Errorf("unexpected model type")
	}

	finalCfg := wm.Config()
	wd, _ := os.Getwd()

	projectDir := wd
	if finalCfg.IsNewProject() && finalCfg.ProjectName != "" {
		projectDir = filepath.Join(wd, finalCfg.ProjectName)
	}

	runner := &executor.RealRunner{}
	ex := executor.New(finalCfg, runner, projectDir)
	ex.SetProjectName(finalCfg.ProjectName)
	tasks := ex.Tasks()
	wizardTasks := executor.ToWizardTasks(tasks)

	execStep := wsteps.NewExecuteStep(wizardTasks)
	execModel := wizard.NewExecuteModel(execStep, tasks)
	ep := tea.NewProgram(execModel, tea.WithAltScreen())
	_, err = ep.Run()
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)
}
