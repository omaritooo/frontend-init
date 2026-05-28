package wizard

import (
	"github.com/omaritooo/frontend-init/config"
	"github.com/omaritooo/frontend-init/wizard/steps"
)

// BuildSteps builds the step slice for the wizard.
// If a preset has been applied (cfg.Preset != "custom"), skips to confirm.
// Otherwise returns the full custom wizard flow.
func BuildSteps(cfg *config.ProjectConfig) []steps.Step {
	if cfg.Preset != "custom" && cfg.Preset != "" {
		return []steps.Step{
			steps.NewSelectStep("preset", presetChoices()),
			steps.NewConfirmStep(cfg),
		}
	}
	return BuildInitialSteps(cfg)
}

// BuildInitialSteps returns the full ordered step slice for custom setup.
func BuildInitialSteps(cfg *config.ProjectConfig) []steps.Step {
	s := []steps.Step{
		steps.NewSelectStep("mode", []string{"new", "existing"}),
		steps.NewSelectStep("preset", append(presetChoices(), "Custom")),
		steps.NewSelectStep("package manager", []string{"npm", "pnpm", "yarn", "bun"}),
		steps.NewSelectStep("framework", []string{"react", "vue", "svelte", "angular", "astro"}),
	}
	s = append(s, variantSteps(cfg)...)
	s = append(s,
		steps.NewSelectStep("typescript", []string{"yes", "no"}),
		steps.NewSelectStep("linting", []string{"eslint-prettier", "biome", "oxlint", "none"}),
		steps.NewSelectStep("ui library", UILibraryOptions(cfg)),
		steps.NewMultiSelectStep("testing", TestingOptions(cfg)),
		steps.NewMultiSelectStep("tooling", ToolingOptions(cfg)),
		steps.NewConfirmStep(cfg),
	)
	return s
}

func variantSteps(cfg *config.ProjectConfig) []steps.Step {
	switch cfg.Framework {
	case "react":
		return []steps.Step{steps.NewSelectStep("variant", []string{"vite", "nextjs"})}
	case "vue":
		return []steps.Step{steps.NewSelectStep("variant", []string{"vite", "nuxt"})}
	case "svelte":
		return []steps.Step{steps.NewSelectStep("variant", []string{"vite", "sveltekit"})}
	case "angular":
		return []steps.Step{steps.NewSelectStep("variant", []string{"angular-cli", "analog"})}
	case "astro":
		return []steps.Step{steps.NewSelectStep("variant", []string{"static", "ssr"})}
	}
	return nil
}

// UILibraryOptions returns UI library choices filtered by framework.
func UILibraryOptions(cfg *config.ProjectConfig) []string {
	base := []string{"none", "tailwind-only"}
	switch cfg.Framework {
	case "react":
		return append(base, "shadcn", "mui", "mantine", "chakra", "antd", "primereact", "daisyui")
	case "vue":
		return append(base, "vuetify", "primevue", "naive-ui", "daisyui")
	case "angular":
		return append(base, "angular-material", "primeng", "ng-zorro")
	case "svelte":
		return append(base, "shadcn-svelte", "skeleton-ui", "daisyui")
	case "astro":
		return append(base, "daisyui", "shadcn")
	}
	return base
}

// TestingOptions returns testing choices filtered by framework.
func TestingOptions(cfg *config.ProjectConfig) []string {
	switch cfg.Framework {
	case "angular":
		return []string{"jest", "playwright", "cypress"}
	default:
		return []string{"vitest", "jest", "testing-library", "playwright", "cypress", "storybook"}
	}
}

// ToolingOptions returns tooling choices filtered by framework and variant.
func ToolingOptions(cfg *config.ProjectConfig) []string {
	switch cfg.Framework {
	case "react":
		base := []string{"tanstack-query", "zustand", "jotai", "redux-toolkit", "rhf-zod", "zod", "axios", "i18next"}
		if cfg.Variant == "nextjs" {
			return append(base, "trpc")
		}
		return append([]string{"tanstack-router", "react-router-v7"}, base...)
	case "vue":
		return []string{"pinia", "tanstack-query", "veevalidate-zod", "axios", "vue-i18n"}
	case "angular":
		return []string{"ngrx-signals", "axios"}
	case "svelte":
		return []string{"tanstack-query", "superforms-zod"}
	case "astro":
		return []string{"nanostores", "zod"}
	}
	return nil
}

func shadcnThemes() []string {
	return []string{"zinc", "slate", "gray", "neutral", "stone", "red", "rose", "orange", "green", "blue", "violet", "yellow"}
}

func presetChoices() []string {
	presets := config.AllPresets()
	names := make([]string, len(presets))
	for i, p := range presets {
		names[i] = p.Name
	}
	return names
}
