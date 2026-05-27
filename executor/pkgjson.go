package executor

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// MergeScripts merges new script entries into the project's package.json.
func MergeScripts(projectDir string, scripts map[string]string) error {
	if len(scripts) == 0 {
		return nil
	}
	pkgPath := filepath.Join(projectDir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return err
	}
	var pkg map[string]any
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}
	existing, ok := pkg["scripts"].(map[string]any)
	if !ok {
		existing = make(map[string]any)
	}
	for k, v := range scripts {
		existing[k] = v
	}
	pkg["scripts"] = existing
	out, err := json.MarshalIndent(pkg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(pkgPath, out, 0644)
}
