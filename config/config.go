package config

type ProjectConfig struct {
	Mode           string   // "new" | "existing"
	ProjectName    string   // only set for new projects
	Preset         string   // preset name or "custom"
	PackageManager string   // npm | pnpm | yarn | bun
	Framework      string   // react | vue | svelte | angular | astro
	Variant        string   // vite | nextjs | nuxt | sveltekit | analog | static | ssr
	TypeScript     bool
	Linting        string   // eslint-prettier | biome | oxlint | none
	UILibrary      string
	ShadcnTheme    string   // base color for shadcn init: zinc | slate | gray | neutral | stone | red | rose | orange | green | blue | violet | yellow
	Testing        []string
	Tooling        []string
}

func New() *ProjectConfig {
	return &ProjectConfig{
		PackageManager: "npm",
		TypeScript:     true,
		Preset:         "custom",
	}
}

func (c *ProjectConfig) IsNewProject() bool {
	return c.Mode == "new"
}
