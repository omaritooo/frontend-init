package executor

import (
	"os"
	"path/filepath"
)

// DetectPackageManager infers the package manager from lockfiles in dir.
func DetectPackageManager(dir string) string {
	lockfiles := map[string]string{
		"pnpm-lock.yaml":    "pnpm",
		"yarn.lock":         "yarn",
		"bun.lockb":         "bun",
		"package-lock.json": "npm",
	}
	for file, pm := range lockfiles {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			return pm
		}
	}
	return "npm"
}
