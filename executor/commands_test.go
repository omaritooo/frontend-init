package executor_test

import (
	"testing"

	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/executor"
	"github.com/stretchr/testify/assert"
)

func TestScaffoldCmd_ReactVite(t *testing.T) {
	cfg := &config.ProjectConfig{
		Framework: "react", Variant: "vite",
		PackageManager: "npm", TypeScript: true,
	}
	cmd := executor.ScaffoldCommand(cfg, "my-app")
	assert.Equal(t, "npm", cmd[0])
	assert.Contains(t, cmd, "create")
	assert.Contains(t, cmd, "vite@latest")
}

func TestScaffoldCmd_NextJS(t *testing.T) {
	cfg := &config.ProjectConfig{
		Framework: "react", Variant: "nextjs",
		PackageManager: "pnpm",
	}
	cmd := executor.ScaffoldCommand(cfg, "my-app")
	assert.Equal(t, "pnpm", cmd[0])
	assert.Contains(t, cmd, "dlx")
	assert.Contains(t, cmd, "create-next-app@latest")
}

func TestScaffoldCmd_Angular(t *testing.T) {
	cfg := &config.ProjectConfig{
		Framework: "angular", Variant: "angular-cli",
		PackageManager: "npm",
	}
	cmd := executor.ScaffoldCommand(cfg, "my-app")
	assert.Contains(t, cmd, "ng")
	assert.Contains(t, cmd, "new")
}
