package executor_test

import (
	"testing"

	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/executor"
	"github.com/stretchr/testify/assert"
)

type mockRunner struct {
	ran [][]string
}

func (m *mockRunner) Run(dir, name string, args ...string) error {
	m.ran = append(m.ran, append([]string{name}, args...))
	return nil
}

func TestExecutor_RunsScaffoldForNewProject(t *testing.T) {
	cfg := config.New()
	cfg.Mode = "new"
	cfg.Framework = "react"
	cfg.Variant = "vite"
	cfg.Linting = "none"
	cfg.TypeScript = true

	runner := &mockRunner{}
	dir := t.TempDir()
	ex := executor.New(cfg, runner, dir)
	tasks := ex.Tasks()

	assert.True(t, len(tasks) > 0)
	assert.Equal(t, "Scaffold project", tasks[0].Label)
}

func TestExecutor_SkipsScaffoldForExistingProject(t *testing.T) {
	cfg := config.New()
	cfg.Mode = "existing"
	cfg.Framework = "react"
	cfg.Variant = "vite"
	cfg.Linting = "none"

	runner := &mockRunner{}
	dir := t.TempDir()
	ex := executor.New(cfg, runner, dir)
	tasks := ex.Tasks()

	for _, t2 := range tasks {
		assert.NotEqual(t, "Scaffold project", t2.Label)
	}
}
