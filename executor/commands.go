package executor

import (
	"fmt"

	"github.com/omaritooo/frontend-init/config"
)

// CommandRunner abstracts shell execution for testability.
type CommandRunner interface {
	Run(dir string, name string, args ...string) error
}

// ScaffoldCommand returns the shell command to scaffold a new project.
// Returns nil if no scaffold is needed (existing project).
func ScaffoldCommand(cfg *config.ProjectConfig, projectName string) []string {
	pm := cfg.PackageManager

	switch fmt.Sprintf("%s/%s", cfg.Framework, cfg.Variant) {
	case "react/vite":
		return []string{pm, "create", "vite@latest", projectName, "--template",
			boolSelect(cfg.TypeScript, "react-ts", "react")}
	case "react/nextjs":
		return nextjsCmd(pm, projectName, cfg.TypeScript)
	case "vue/vite":
		return []string{pm, "create", "vite@latest", projectName, "--template",
			boolSelect(cfg.TypeScript, "vue-ts", "vue")}
	case "vue/nuxt":
		return []string{pm, "dlx", "nuxi@latest", "init", projectName}
	case "svelte/vite":
		return []string{pm, "create", "vite@latest", projectName, "--template",
			boolSelect(cfg.TypeScript, "svelte-ts", "svelte")}
	case "svelte/sveltekit":
		return []string{pm, "create", "svelte@latest", projectName}
	case "angular/angular-cli":
		return []string{"ng", "new", projectName}
	case "angular/analog":
		return []string{pm, "create", "analog@latest", projectName}
	case "astro/static", "astro/ssr":
		return []string{pm, "create", "astro@latest", projectName}
	}
	return nil
}

func nextjsCmd(pm, name string, ts bool) []string {
	args := []string{pm, "dlx", "create-next-app@latest", name}
	if ts {
		args = append(args, "--typescript")
	}
	return args
}

func boolSelect(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

// InstallCmd returns the command to install packages.
func InstallCmd(pm string, dev bool, packages []string) []string {
	base := []string{pm, "install"}
	if dev {
		base = append(base, "-D")
	}
	return append(base, packages...)
}
