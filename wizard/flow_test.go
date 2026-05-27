package wizard_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/wizard"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

func stepLabels(ss []steps.Step) []string {
	labels := make([]string, len(ss))
	for i, s := range ss {
		labels[i] = s.Label()
	}
	return labels
}

func TestBuildInitialSteps_ContainsCoreScreens(t *testing.T) {
	cfg := config.New()
	all := wizard.BuildInitialSteps(cfg)
	labels := stepLabels(all)
	assert.Contains(t, labels, "mode")
	assert.Contains(t, labels, "framework")
	assert.Contains(t, labels, "linting")
	assert.Contains(t, labels, "confirm")
}

func TestBuildSteps_PresetPathGoesDirectToConfirm(t *testing.T) {
	cfg := config.New()
	p := config.GetPreset("React Minimal")
	p.Apply(cfg)
	all := wizard.BuildSteps(cfg)
	labels := stepLabels(all)
	// preset path should NOT include framework/variant steps (already filled)
	assert.NotContains(t, labels, "framework")
	assert.Contains(t, labels, "confirm")
}

func TestToolingOptions_AngularExcludesTanStackRouter(t *testing.T) {
	cfg := config.New()
	cfg.Framework = "angular"
	opts := wizard.ToolingOptions(cfg)
	assert.NotContains(t, opts, "tanstack-router")
}

func TestToolingOptions_ReactViteIncludesTanStackRouter(t *testing.T) {
	cfg := config.New()
	cfg.Framework = "react"
	cfg.Variant = "vite"
	opts := wizard.ToolingOptions(cfg)
	assert.Contains(t, opts, "tanstack-router")
}

func TestToolingOptions_NextJsIncludesTRPC(t *testing.T) {
	cfg := config.New()
	cfg.Framework = "react"
	cfg.Variant = "nextjs"
	opts := wizard.ToolingOptions(cfg)
	assert.Contains(t, opts, "trpc")
	assert.NotContains(t, opts, "tanstack-router")
}

func TestUILibraryOptions_ReactHasShadcn(t *testing.T) {
	cfg := config.New()
	cfg.Framework = "react"
	opts := wizard.UILibraryOptions(cfg)
	assert.Contains(t, opts, "shadcn")
}

func TestUILibraryOptions_AngularHasAngularMaterial(t *testing.T) {
	cfg := config.New()
	cfg.Framework = "angular"
	opts := wizard.UILibraryOptions(cfg)
	assert.Contains(t, opts, "angular-material")
	assert.NotContains(t, opts, "shadcn")
}
