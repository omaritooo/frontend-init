package executor_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/omaritooo/frontend-init/executor"
)

func TestPatchFile_InsertAfter(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "src/main.tsx")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	original := "import React from 'react'\nReactDOM.createRoot(document.getElementById('root')!)"
	_ = os.WriteFile(p, []byte(original), 0644)

	patch := executor.FilePatch{
		Path:   "src/main.tsx",
		Find:   "import React from 'react'",
		Insert: "import { QueryClient } from '@tanstack/react-query'\n",
		Mode:   executor.PatchInsertAfter,
	}
	err := executor.ApplyPatch(dir, patch)
	assert.NoError(t, err)
	data, _ := os.ReadFile(p)
	assert.Contains(t, string(data), "QueryClient")
	assert.Contains(t, string(data), "import React")
}

func TestPatchFile_Append(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "src/index.css")
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	_ = os.WriteFile(p, []byte("body { margin: 0; }"), 0644)

	patch := executor.FilePatch{
		Path:   "src/index.css",
		Insert: `@import "tailwindcss";`,
		Mode:   executor.PatchAppend,
	}
	err := executor.ApplyPatch(dir, patch)
	assert.NoError(t, err)
	data, _ := os.ReadFile(p)
	assert.Contains(t, string(data), `@import "tailwindcss"`)
}

func TestPatchFile_SkipsIfFileNotFound(t *testing.T) {
	dir := t.TempDir()
	patch := executor.FilePatch{
		Path:   "does/not/exist.ts",
		Insert: "foo",
		Mode:   executor.PatchAppend,
	}
	err := executor.ApplyPatch(dir, patch)
	assert.NoError(t, err)
}
