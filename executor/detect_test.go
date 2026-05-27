package executor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/omaritooo/frontend-init/executor"
)

func TestDetectPackageManager_FromLockfile(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"pnpm-lock.yaml", "pnpm"},
		{"yarn.lock", "yarn"},
		{"bun.lockb", "bun"},
		{"package-lock.json", "npm"},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			dir := t.TempDir()
			_ = os.WriteFile(filepath.Join(dir, tt.file), []byte(""), 0644)
			got := executor.DetectPackageManager(dir)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDetectPackageManager_DefaultsToNpm(t *testing.T) {
	dir := t.TempDir()
	got := executor.DetectPackageManager(dir)
	assert.Equal(t, "npm", got)
}
