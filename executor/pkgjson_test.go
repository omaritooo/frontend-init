package executor_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/executor"
)

func TestMergeScripts_AddsNewScripts(t *testing.T) {
	dir := t.TempDir()
	pkg := map[string]any{
		"name":    "my-app",
		"scripts": map[string]any{"dev": "vite"},
	}
	data, _ := json.Marshal(pkg)
	_ = os.WriteFile(filepath.Join(dir, "package.json"), data, 0644)

	err := executor.MergeScripts(dir, map[string]string{
		"lint":   "eslint .",
		"format": "prettier --write .",
	})
	assert.NoError(t, err)

	result, _ := os.ReadFile(filepath.Join(dir, "package.json"))
	var out map[string]any
	_ = json.Unmarshal(result, &out)
	scripts := out["scripts"].(map[string]any)
	assert.Equal(t, "eslint .", scripts["lint"])
	assert.Equal(t, "vite", scripts["dev"]) // existing preserved
}
