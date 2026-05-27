package config_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/config"
)

func TestPresets_AllDefined(t *testing.T) {
	presets := config.AllPresets()
	assert.NotEmpty(t, presets)
	names := make([]string, len(presets))
	for i, p := range presets {
		names[i] = p.Name
	}
	assert.Contains(t, names, "React Minimal")
	assert.Contains(t, names, "T3 Stack")
	assert.Contains(t, names, "Angular Enterprise")
	assert.Contains(t, names, "Astro Islands")
}

func TestPresets_ApplyToConfig(t *testing.T) {
	cfg := config.New()
	p := config.GetPreset("React Minimal")
	assert.NotNil(t, p)
	p.Apply(cfg)
	assert.Equal(t, "react", cfg.Framework)
	assert.Equal(t, "vite", cfg.Variant)
	assert.Equal(t, "eslint-prettier", cfg.Linting)
	assert.Contains(t, cfg.Testing, "vitest")
}
