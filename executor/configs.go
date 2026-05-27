package executor

import (
	"os"
	"path/filepath"
)

type PatchMode int

const (
	PatchAppend      PatchMode = iota // append Insert to end of file
	PatchInsertAfter                  // insert Insert after first occurrence of Find
	PatchReplace                      // replace all occurrences of Find with Insert
)

type ConfigFile struct {
	Path    string
	Content string
}

type FilePatch struct {
	Path   string
	Find   string
	Insert string
	Mode   PatchMode
}

type ToolSetup struct {
	Name            string
	Packages        []string
	DevPackages     []string
	ConfigFiles     []ConfigFile
	FilePatches     []FilePatch
	PostInstallCmds []string
	Scripts         map[string]string
}

// GetToolSetup returns the ToolSetup for the given tool key scoped to a framework.
// Returns nil if the tool is not in the catalog.
func GetToolSetup(tool, framework string) *ToolSetup {
	catalog := buildCatalog(framework)
	if s, ok := catalog[tool]; ok {
		return &s
	}
	return nil
}

func buildCatalog(framework string) map[string]ToolSetup {
	mainEntry := mainEntryFile(framework)
	cssEntry := cssEntryFile(framework)

	return map[string]ToolSetup{
		"tailwind": {
			Name:        "Tailwind CSS",
			DevPackages: []string{"tailwindcss", "@tailwindcss/vite"},
			ConfigFiles: []ConfigFile{
				{Path: "tailwind.config.ts", Content: tailwindConfig},
			},
			FilePatches: []FilePatch{
				{Path: cssEntry, Insert: `@import "tailwindcss";`, Mode: PatchAppend},
			},
			Scripts: map[string]string{},
		},
		"shadcn": {
			Name:            "shadcn/ui",
			DevPackages:     []string{"shadcn"},
			PostInstallCmds: []string{"npx shadcn init"},
		},
		"eslint-prettier": {
			Name:        "ESLint + Prettier",
			DevPackages: []string{"eslint", "prettier", "eslint-config-prettier", "@eslint/js"},
			ConfigFiles: []ConfigFile{
				{Path: "eslint.config.js", Content: eslintConfig},
				{Path: ".prettierrc", Content: prettierConfig},
				{Path: ".prettierignore", Content: prettierIgnore},
			},
			Scripts: map[string]string{
				"lint":         "eslint .",
				"lint:fix":     "eslint . --fix",
				"format":       "prettier --write .",
				"format:check": "prettier --check .",
			},
		},
		"biome": {
			Name:            "Biome",
			DevPackages:     []string{"@biomejs/biome"},
			PostInstallCmds: []string{"npx biome init"},
			Scripts: map[string]string{
				"lint":   "biome lint .",
				"format": "biome format --write .",
				"check":  "biome check .",
			},
		},
		"oxlint": {
			Name:        "Oxlint + Prettier",
			DevPackages: []string{"oxlint", "prettier"},
			ConfigFiles: []ConfigFile{
				{Path: ".prettierrc", Content: prettierConfig},
			},
			Scripts: map[string]string{
				"lint":   "oxlint .",
				"format": "prettier --write .",
			},
		},
		"vitest": {
			Name:        "Vitest",
			DevPackages: []string{"vitest", "@vitest/ui"},
			ConfigFiles: []ConfigFile{
				{Path: "vitest.config.ts", Content: vitestConfig},
			},
			Scripts: map[string]string{
				"test":     "vitest",
				"test:ui":  "vitest --ui",
				"coverage": "vitest run --coverage",
			},
		},
		"jest": {
			Name:        "Jest",
			DevPackages: []string{"jest", "ts-jest", "@types/jest"},
			ConfigFiles: []ConfigFile{
				{Path: "jest.config.ts", Content: jestConfig},
			},
			Scripts: map[string]string{"test": "jest"},
		},
		"testing-library": {
			Name:        "Testing Library",
			DevPackages: []string{"@testing-library/react", "@testing-library/user-event", "@testing-library/jest-dom"},
		},
		"playwright": {
			Name:            "Playwright",
			DevPackages:     []string{"@playwright/test"},
			ConfigFiles:     []ConfigFile{{Path: "playwright.config.ts", Content: playwrightConfig}},
			PostInstallCmds: []string{"npx playwright install"},
			Scripts:         map[string]string{"e2e": "playwright test", "e2e:ui": "playwright test --ui"},
		},
		"cypress": {
			Name:        "Cypress",
			DevPackages: []string{"cypress"},
			Scripts:     map[string]string{"e2e": "cypress run", "e2e:open": "cypress open"},
		},
		"storybook": {
			Name:            "Storybook",
			DevPackages:     []string{},
			PostInstallCmds: []string{"npx storybook@latest init"},
			Scripts:         map[string]string{"storybook": "storybook dev -p 6006", "build-storybook": "storybook build"},
		},
		"tanstack-query": func() ToolSetup {
			switch framework {
			case "vue", "nuxt":
				return ToolSetup{
					Name:     "TanStack Query",
					Packages: []string{"@tanstack/vue-query"},
				}
			case "svelte", "sveltekit":
				return ToolSetup{
					Name:     "TanStack Query",
					Packages: []string{"@tanstack/svelte-query"},
				}
			default: // react, nextjs, astro
				return ToolSetup{
					Name:     "TanStack Query",
					Packages: []string{"@tanstack/react-query"},
					FilePatches: []FilePatch{
						{
							Path:   mainEntry,
							Find:   "ReactDOM.createRoot",
							Insert: "import { QueryClient, QueryClientProvider } from '@tanstack/react-query'\nconst queryClient = new QueryClient()\n",
							Mode:   PatchInsertAfter,
						},
					},
				}
			}
		}(),
		"tanstack-router": {
			Name:        "TanStack Router",
			Packages:    []string{"@tanstack/react-router"},
			DevPackages: []string{"@tanstack/router-plugin"},
		},
		"react-router-v7": {
			Name:     "React Router v7",
			Packages: []string{"react-router"},
		},
		"zustand": {
			Name:     "Zustand",
			Packages: []string{"zustand"},
		},
		"jotai": {
			Name:     "Jotai",
			Packages: []string{"jotai"},
		},
		"redux-toolkit": {
			Name:     "Redux Toolkit",
			Packages: []string{"@reduxjs/toolkit", "react-redux"},
		},
		"rhf-zod": {
			Name:     "React Hook Form + Zod",
			Packages: []string{"react-hook-form", "zod", "@hookform/resolvers"},
		},
		"zod": {
			Name:     "Zod",
			Packages: []string{"zod"},
		},
		"trpc": {
			Name:     "tRPC",
			Packages: []string{"@trpc/server", "@trpc/client", "@trpc/react-query"},
		},
		"axios": {
			Name:     "Axios",
			Packages: []string{"axios"},
		},
		"i18next": {
			Name:     "i18next",
			Packages: []string{"i18next", "react-i18next"},
		},
		"pinia": {
			Name:     "Pinia",
			Packages: []string{"pinia"},
			FilePatches: []FilePatch{
				{
					Path:   mainEntry,
					Find:   "app.mount",
					Insert: "import { createPinia } from 'pinia'\napp.use(createPinia())\n",
					Mode:   PatchInsertAfter,
				},
			},
		},
		"veevalidate-zod": {
			Name:     "VeeValidate + Zod",
			Packages: []string{"vee-validate", "@vee-validate/zod", "zod"},
		},
		"vue-i18n": {
			Name:     "vue-i18n",
			Packages: []string{"vue-i18n"},
		},
		"ngrx-signals": {
			Name:     "NgRx Signal Store",
			Packages: []string{"@ngrx/signals"},
		},
		"superforms-zod": {
			Name:     "Superforms + Zod",
			Packages: []string{"sveltekit-superforms", "zod"},
		},
		"nanostores": {
			Name:     "Nanostores",
			Packages: []string{"nanostores"},
		},
		"vuetify": {
			Name:        "Vuetify",
			Packages:    []string{"vuetify"},
			DevPackages: []string{"vite-plugin-vuetify"},
		},
		"primevue": {
			Name:     "PrimeVue",
			Packages: []string{"primevue", "@primevue/themes"},
		},
		"naive-ui": {
			Name:     "Naive UI",
			Packages: []string{"naive-ui"},
		},
		"daisyui": {
			Name:        "DaisyUI",
			DevPackages: []string{"daisyui"},
		},
		"mantine": {
			Name:     "Mantine",
			Packages: []string{"@mantine/core", "@mantine/hooks"},
		},
		"chakra": {
			Name:     "Chakra UI",
			Packages: []string{"@chakra-ui/react"},
		},
		"mui": {
			Name:     "MUI",
			Packages: []string{"@mui/material", "@emotion/react", "@emotion/styled"},
		},
		"antd": {
			Name:     "Ant Design",
			Packages: []string{"antd"},
		},
		"primereact": {
			Name:     "PrimeReact",
			Packages: []string{"primereact", "primeicons"},
		},
		"angular-material": {
			Name:            "Angular Material",
			PostInstallCmds: []string{"ng add @angular/material"},
		},
		"primeng": {
			Name:     "PrimeNG",
			Packages: []string{"primeng", "primeicons"},
		},
		"ng-zorro": {
			Name:            "NG-Zorro",
			PostInstallCmds: []string{"ng add ng-zorro-antd"},
		},
		"skeleton-ui": {
			Name:     "Skeleton UI",
			Packages: []string{"@skeletonlabs/skeleton"},
		},
		"shadcn-svelte": {
			Name:            "shadcn-svelte",
			DevPackages:     []string{},
			PostInstallCmds: []string{"npx shadcn-svelte@latest init"},
		},
	}
}

func mainEntryFile(framework string) string {
	switch framework {
	case "vue", "angular":
		return "src/main.ts"
	default:
		return "src/main.tsx"
	}
}

func cssEntryFile(_ string) string {
	return "src/index.css"
}

const tailwindConfig = `import type { Config } from 'tailwindcss'
export default {
  content: ['./index.html', './src/**/*.{ts,tsx,vue,svelte}'],
  theme: { extend: {} },
  plugins: [],
} satisfies Config
`

const eslintConfig = `import js from '@eslint/js'
export default [
  js.configs.recommended,
  { rules: { 'no-unused-vars': 'warn' } },
]
`

const prettierConfig = `{
  "semi": false,
  "singleQuote": true,
  "tabWidth": 2,
  "trailingComma": "es5"
}
`

const prettierIgnore = `node_modules
dist
.next
.nuxt
`

const vitestConfig = `import { defineConfig } from 'vitest/config'
export default defineConfig({
  test: { environment: 'jsdom' },
})
`

const jestConfig = `export default {
  preset: 'ts-jest',
  testEnvironment: 'node',
}
`

const playwrightConfig = `import { defineConfig } from '@playwright/test'
export default defineConfig({
  testDir: './e2e',
  use: { baseURL: 'http://localhost:5173' },
})
`

// WriteConfigFiles writes all config files relative to projectDir.
func WriteConfigFiles(projectDir string, files []ConfigFile) error {
	for _, f := range files {
		fullPath := filepath.Join(projectDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(fullPath, []byte(f.Content), 0644); err != nil {
			return err
		}
	}
	return nil
}
