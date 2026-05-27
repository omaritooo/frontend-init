package executor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/executor"
)

func TestWriteConfigFiles_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	files := []executor.ConfigFile{
		{Path: "eslint.config.js", Content: "export default []"},
		{Path: ".prettierrc", Content: `{"semi": false}`},
	}
	err := executor.WriteConfigFiles(dir, files)
	assert.NoError(t, err)
	for _, f := range files {
		data, err := os.ReadFile(filepath.Join(dir, f.Path))
		assert.NoError(t, err)
		assert.Equal(t, f.Content, string(data))
	}
}

func TestWriteConfigFiles_CreatesSubdirs(t *testing.T) {
	dir := t.TempDir()
	files := []executor.ConfigFile{
		{Path: "src/config/env.ts", Content: "export const env = {}"},
	}
	err := executor.WriteConfigFiles(dir, files)
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(dir, "src/config/env.ts"))
	assert.NoError(t, err)
}
