package executor_test

import (
	"testing"

	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/executor"
	"github.com/stretchr/testify/assert"
)

func reactCfg() *config.ProjectConfig {
	return &config.ProjectConfig{Framework: "react"}
}

func TestToolCatalog_TailwindHasFilePatch(t *testing.T) {
	setup := executor.GetToolSetup("tailwind", reactCfg())
	assert.NotNil(t, setup)
	assert.NotEmpty(t, setup.DevPackages)
	// Tailwind v4 uses CSS imports and vite plugin — no tailwind.config.ts needed
	assert.NotEmpty(t, setup.FilePatches)
}

func TestToolCatalog_ShadcnDefaultUsesZinc(t *testing.T) {
	setup := executor.GetToolSetup("shadcn", reactCfg())
	assert.NotNil(t, setup)
	assert.Contains(t, setup.PostInstallCmds, "npx shadcn@latest init -d")
}

func TestToolCatalog_ShadcnCustomThemePassesBaseColor(t *testing.T) {
	cfg := &config.ProjectConfig{Framework: "react", ShadcnTheme: "slate"}
	setup := executor.GetToolSetup("shadcn", cfg)
	assert.NotNil(t, setup)
	assert.Contains(t, setup.PostInstallCmds, "npx shadcn@latest init -d --base-color slate")
}

func TestToolCatalog_TanStackQueryPatchesMainTsx(t *testing.T) {
	setup := executor.GetToolSetup("tanstack-query", reactCfg())
	assert.NotNil(t, setup)
	found := false
	for _, p := range setup.FilePatches {
		if p.Path == "src/main.tsx" {
			found = true
		}
	}
	assert.True(t, found, "tanstack-query should patch src/main.tsx")
}

func TestToolCatalog_PlaywrightHasPostInstallCmd(t *testing.T) {
	setup := executor.GetToolSetup("playwright", reactCfg())
	assert.NotNil(t, setup)
	assert.Contains(t, setup.PostInstallCmds, "npx playwright install")
}

func TestToolCatalog_UnknownToolReturnsNil(t *testing.T) {
	setup := executor.GetToolSetup("does-not-exist", reactCfg())
	assert.Nil(t, setup)
}
