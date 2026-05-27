package config_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/config"
)

func TestProjectConfig_Defaults(t *testing.T) {
	cfg := config.New()
	assert.Equal(t, "npm", cfg.PackageManager)
	assert.True(t, cfg.TypeScript)
	assert.Equal(t, "custom", cfg.Preset)
}

func TestProjectConfig_IsNewProject(t *testing.T) {
	cfg := config.New()
	cfg.Mode = "new"
	assert.True(t, cfg.IsNewProject())
	cfg.Mode = "existing"
	assert.False(t, cfg.IsNewProject())
}
