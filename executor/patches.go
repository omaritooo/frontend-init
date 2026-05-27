package executor

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ApplyPatch applies a single FilePatch to a file in projectDir.
// Missing files are silently skipped.
func ApplyPatch(projectDir string, patch FilePatch) error {
	fullPath := filepath.Join(projectDir, patch.Path)
	data, err := os.ReadFile(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	content := string(data)

	switch patch.Mode {
	case PatchAppend:
		content = content + "\n" + patch.Insert
	case PatchInsertAfter:
		idx := strings.Index(content, patch.Find)
		if idx == -1 {
			content = patch.Insert + "\n" + content
		} else {
			insertAt := idx + len(patch.Find)
			content = content[:insertAt] + "\n" + patch.Insert + content[insertAt:]
		}
	case PatchReplace:
		content = strings.ReplaceAll(content, patch.Find, patch.Insert)
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}

// ApplyPatches applies all patches sequentially.
func ApplyPatches(projectDir string, patches []FilePatch) error {
	for _, p := range patches {
		if err := ApplyPatch(projectDir, p); err != nil {
			return err
		}
	}
	return nil
}
