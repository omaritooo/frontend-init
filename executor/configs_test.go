package executor_test

import (
	"testing"

	"github.com/omaritooo/frontend-init/executor"
	"github.com/stretchr/testify/assert"
)

func TestToolCatalog_TailwindHasFilePatch(t *testing.T) {
	setup := executor.GetToolSetup("tailwind", "react")
	assert.NotNil(t, setup)
	assert.NotEmpty(t, setup.DevPackages)
	// Tailwind v4 uses CSS imports and vite plugin — no tailwind.config.ts needed
	assert.NotEmpty(t, setup.FilePatches)
}

func TestToolCatalog_ShadcnHasPostInstallCmd(t *testing.T) {
	setup := executor.GetToolSetup("shadcn", "react")
	assert.NotNil(t, setup)
	assert.Contains(t, setup.PostInstallCmds, "npx shadcn@latest init --yes")
}

func TestToolCatalog_TanStackQueryPatchesMainTsx(t *testing.T) {
	setup := executor.GetToolSetup("tanstack-query", "react")
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
	setup := executor.GetToolSetup("playwright", "react")
	assert.NotNil(t, setup)
	assert.Contains(t, setup.PostInstallCmds, "npx playwright install")
}

func TestToolCatalog_UnknownToolReturnsNil(t *testing.T) {
	setup := executor.GetToolSetup("does-not-exist", "react")
	assert.Nil(t, setup)
}
