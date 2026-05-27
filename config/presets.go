package config

type Preset struct {
	Name      string
	Framework string
	Variant   string
	Linting   string
	UILibrary string
	Testing   []string
	Tooling   []string
}

func (p *Preset) Apply(cfg *ProjectConfig) {
	cfg.Preset     = p.Name
	cfg.Framework  = p.Framework
	cfg.Variant    = p.Variant
	cfg.Linting    = p.Linting
	cfg.UILibrary  = p.UILibrary
	cfg.Testing    = append([]string{}, p.Testing...)
	cfg.Tooling    = append([]string{}, p.Tooling...)
	cfg.TypeScript = true
}

func AllPresets() []Preset {
	return []Preset{
		{
			Name: "React Minimal", Framework: "react", Variant: "vite",
			Linting: "eslint-prettier", Testing: []string{"vitest"},
		},
		{
			Name: "React Full SPA", Framework: "react", Variant: "vite",
			Linting: "eslint-prettier", UILibrary: "shadcn",
			Testing: []string{"vitest", "testing-library", "playwright"},
			Tooling: []string{"tanstack-query", "tanstack-router", "zustand", "rhf-zod"},
		},
		{
			Name: "Next.js Standard", Framework: "react", Variant: "nextjs",
			Linting: "eslint-prettier", UILibrary: "shadcn",
			Testing: []string{"vitest"},
			Tooling: []string{"tanstack-query", "zod"},
		},
		{
			Name: "T3 Stack", Framework: "react", Variant: "nextjs",
			Linting: "eslint-prettier", UILibrary: "shadcn",
			Testing: []string{"vitest"},
			Tooling: []string{"trpc", "tanstack-query", "zod"},
		},
		{
			Name: "Vue Minimal", Framework: "vue", Variant: "vite",
			Linting: "eslint-prettier", Testing: []string{"vitest"},
		},
		{
			Name: "Vue Full SPA", Framework: "vue", Variant: "vite",
			Linting: "eslint-prettier", UILibrary: "daisyui",
			Testing: []string{"vitest"},
			Tooling: []string{"pinia", "tanstack-query", "veevalidate-zod"},
		},
		{
			Name: "Nuxt Standard", Framework: "vue", Variant: "nuxt",
			Linting: "eslint-prettier", Testing: []string{"vitest"},
		},
		{
			Name: "Angular Minimal", Framework: "angular", Variant: "angular-cli",
			Linting: "eslint-prettier", Testing: []string{"jest"},
		},
		{
			Name: "Angular Enterprise", Framework: "angular", Variant: "angular-cli",
			Linting: "eslint-prettier", UILibrary: "angular-material",
			Testing: []string{"jest", "playwright"},
			Tooling: []string{"ngrx-signals"},
		},
		{
			Name: "Astro Content Site", Framework: "astro", Variant: "static",
			Linting: "eslint-prettier", UILibrary: "daisyui",
			Testing: []string{"vitest"},
		},
		{
			Name: "Astro Islands", Framework: "astro", Variant: "ssr",
			Linting: "eslint-prettier", UILibrary: "shadcn",
			Testing: []string{"vitest"},
			Tooling: []string{"tanstack-query"},
		},
	}
}

func GetPreset(name string) *Preset {
	for _, p := range AllPresets() {
		if p.Name == name {
			cp := p
			return &cp
		}
	}
	return nil
}
