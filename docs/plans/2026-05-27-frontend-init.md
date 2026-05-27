# frontend-init Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go CLI that walks users through a Bubbletea TUI wizard to scaffold or configure frontend projects with framework-filtered tooling, named presets, and automated multi-step tool setup.

**Architecture:** A Cobra `init` command launches a Bubbletea app whose root model holds a dynamic `[]Step` slice. Framework-conditional steps are inserted at runtime after the framework choice resolves. After confirmation, a sequential executor runs each selected tool's full setup sequence: install packages → write config files → patch existing files → run post-install commands.

**Tech Stack:** Go 1.22+, Cobra v1.10, Bubbletea v1.3, Lipgloss v1.1, Testify v1.10

**Design doc:** `docs/plans/2026-05-27-frontend-init-design.md`

---

## Task 1: Project skeleton

**Files:**
- Create: `config/config.go`
- Create: `config/presets.go`
- Create: `executor/executor.go`
- Create: `executor/commands.go`
- Create: `executor/configs.go`
- Create: `executor/patches.go`
- Create: `wizard/model.go`
- Create: `wizard/steps/step.go`
- Create: `wizard/steps/select.go`
- Create: `wizard/steps/multiselect.go`
- Create: `wizard/steps/confirm.go`
- Create: `wizard/steps/execute.go`
- Create: `wizard/steps.go`

**Step 1: Create all directories**

```bash
mkdir -p config executor wizard/steps
```

**Step 2: Create stub files with package declarations**

Each file needs only its `package` declaration for now. Example:
- `config/config.go` → `package config`
- `executor/executor.go` → `package executor`
- `wizard/model.go` → `package wizard`
- `wizard/steps/step.go` → `package steps`
- `wizard/steps.go` → `package wizard`

**Step 3: Add testify**

```bash
go get github.com/stretchr/testify@latest
go mod tidy
```

**Step 4: Verify build**

```bash
go build ./...
```
Expected: no errors.

**Step 5: Commit**

```bash
git init
git add .
git commit -m "chore: project skeleton with package stubs"
```

---

## Task 2: ProjectConfig struct

**Files:**
- Modify: `config/config.go`
- Create: `config/config_test.go`

**Step 1: Write the failing test**

`config/config_test.go`:
```go
package config_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
)

func TestProjectConfig_Defaults(t *testing.T) {
    cfg := config.New()
    assert.Equal(t, "npm", cfg.PackageManager)
    assert.True(t, cfg.TypeScript)
    assert.Equal(t, "custom", cfg.Preset)
}

func TestProjectConfig_IsNewProject(t *testing.T) {
    cfg := config.New()
    cfg.Mode = "new"
    assert.True(t, cfg.IsNewProject())
    cfg.Mode = "existing"
    assert.False(t, cfg.IsNewProject())
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./config/... -v
```
Expected: FAIL — `config.New` undefined.

**Step 3: Implement ProjectConfig**

`config/config.go`:
```go
package config

type ProjectConfig struct {
    Mode           string   // "new" | "existing"
    Preset         string   // preset name or "custom"
    PackageManager string   // npm | pnpm | yarn | bun
    Framework      string   // react | vue | svelte | angular | astro
    Variant        string   // vite | nextjs | nuxt | sveltekit | analog | static | ssr
    TypeScript     bool
    Linting        string   // eslint-prettier | biome | oxlint | none
    UILibrary      string
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./config/... -v
```
Expected: PASS.

**Step 5: Commit**

```bash
git add config/
git commit -m "feat: ProjectConfig struct with defaults"
```

---

## Task 3: Preset definitions

**Files:**
- Modify: `config/presets.go`
- Create: `config/presets_test.go`

**Step 1: Write the failing test**

`config/presets_test.go`:
```go
package config_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
)

func TestPresets_AllDefined(t *testing.T) {
    presets := config.AllPresets()
    assert.NotEmpty(t, presets)
    names := make([]string, len(presets))
    for i, p := range presets {
        names[i] = p.Name
    }
    assert.Contains(t, names, "React Minimal")
    assert.Contains(t, names, "T3 Stack")
    assert.Contains(t, names, "Angular Enterprise")
    assert.Contains(t, names, "Astro Islands")
}

func TestPresets_ApplyToConfig(t *testing.T) {
    cfg := config.New()
    p := config.GetPreset("React Minimal")
    assert.NotNil(t, p)
    p.Apply(cfg)
    assert.Equal(t, "react", cfg.Framework)
    assert.Equal(t, "vite", cfg.Variant)
    assert.Equal(t, "eslint-prettier", cfg.Linting)
    assert.Contains(t, cfg.Testing, "vitest")
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./config/... -v -run TestPresets
```
Expected: FAIL — `config.AllPresets` undefined.

**Step 3: Implement presets**

`config/presets.go`:
```go
package config

type Preset struct {
    Name           string
    Framework      string
    Variant        string
    Linting        string
    UILibrary      string
    Testing        []string
    Tooling        []string
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
            Name:      "React Minimal",
            Framework: "react", Variant: "vite",
            Linting: "eslint-prettier",
            Testing: []string{"vitest"},
        },
        {
            Name:      "React Full SPA",
            Framework: "react", Variant: "vite",
            Linting: "eslint-prettier", UILibrary: "shadcn",
            Testing: []string{"vitest", "testing-library", "playwright"},
            Tooling: []string{"tanstack-query", "tanstack-router", "zustand", "rhf-zod"},
        },
        {
            Name:      "Next.js Standard",
            Framework: "react", Variant: "nextjs",
            Linting: "eslint-prettier", UILibrary: "shadcn",
            Testing: []string{"vitest"},
            Tooling: []string{"tanstack-query", "zod"},
        },
        {
            Name:      "T3 Stack",
            Framework: "react", Variant: "nextjs",
            Linting: "eslint-prettier", UILibrary: "shadcn",
            Testing: []string{"vitest"},
            Tooling: []string{"trpc", "tanstack-query", "zod"},
        },
        {
            Name:      "Vue Minimal",
            Framework: "vue", Variant: "vite",
            Linting: "eslint-prettier",
            Testing: []string{"vitest"},
        },
        {
            Name:      "Vue Full SPA",
            Framework: "vue", Variant: "vite",
            Linting: "eslint-prettier", UILibrary: "daisyui",
            Testing: []string{"vitest"},
            Tooling: []string{"pinia", "tanstack-query", "veevalidate-zod"},
        },
        {
            Name:      "Nuxt Standard",
            Framework: "vue", Variant: "nuxt",
            Linting: "eslint-prettier",
            Testing: []string{"vitest"},
        },
        {
            Name:      "Angular Minimal",
            Framework: "angular", Variant: "angular-cli",
            Linting: "eslint-prettier",
            Testing: []string{"jest"},
        },
        {
            Name:      "Angular Enterprise",
            Framework: "angular", Variant: "angular-cli",
            Linting: "eslint-prettier", UILibrary: "angular-material",
            Testing: []string{"jest", "playwright"},
            Tooling: []string{"ngrx-signals"},
        },
        {
            Name:      "Astro Content Site",
            Framework: "astro", Variant: "static",
            Linting: "eslint-prettier", UILibrary: "daisyui",
            Testing: []string{"vitest"},
        },
        {
            Name:      "Astro Islands",
            Framework: "astro", Variant: "ssr",
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./config/... -v
```
Expected: PASS.

**Step 5: Commit**

```bash
git add config/
git commit -m "feat: preset definitions for all 11 stacks"
```

---

## Task 4: ToolSetup model and catalog

**Files:**
- Modify: `executor/configs.go`
- Create: `executor/configs_test.go`

**Step 1: Write the failing test**

`executor/configs_test.go`:
```go
package executor_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/executor"
)

func TestToolCatalog_TailwindHasPostInstallPatch(t *testing.T) {
    setup := executor.GetToolSetup("tailwind", "react")
    assert.NotNil(t, setup)
    assert.NotEmpty(t, setup.DevPackages)
    assert.NotEmpty(t, setup.ConfigFiles)
    assert.NotEmpty(t, setup.FilePatches) // must patch index.css
}

func TestToolCatalog_ShadcnHasPostInstallCmd(t *testing.T) {
    setup := executor.GetToolSetup("shadcn", "react")
    assert.NotNil(t, setup)
    assert.Contains(t, setup.PostInstallCmds, "npx shadcn init")
}

func TestToolCatalog_TanStackQueryPatchesMainTsx(t *testing.T) {
    setup := executor.GetToolSetup("tanstack-query", "react")
    assert.NotNil(t, setup)
    found := false
    for _, p := range setup.FilePatches {
        if p.Path == "src/main.tsx" {
            found = true
        }
    }
    assert.True(t, found, "tanstack-query should patch src/main.tsx")
}

func TestToolCatalog_PlaywrightHasPostInstallCmd(t *testing.T) {
    setup := executor.GetToolSetup("playwright", "react")
    assert.NotNil(t, setup)
    assert.Contains(t, setup.PostInstallCmds, "npx playwright install")
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestToolCatalog
```
Expected: FAIL — types not defined.

**Step 3: Define types and catalog**

`executor/configs.go`:
```go
package executor

type PatchMode int

const (
    PatchAppend      PatchMode = iota
    PatchInsertAfter           // insert after first match of Find
    PatchReplace               // replace Find with Insert
)

type ConfigFile struct {
    Path    string // relative to project root
    Content string // file content
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

// GetToolSetup returns the setup for a tool scoped to a framework.
// framework is used for patches that differ by framework (e.g. main.tsx vs main.ts).
func GetToolSetup(tool, framework string) *ToolSetup {
    catalog := buildCatalog(framework)
    if s, ok := catalog[tool]; ok {
        return &s
    }
    return nil
}

func buildCatalog(framework string) map[string]ToolSetup {
    mainEntry := mainEntryFile(framework)
    cssEntry  := cssEntryFile(framework)

    return map[string]ToolSetup{
        "tailwind": {
            Name:        "Tailwind CSS",
            DevPackages: []string{"tailwindcss", "@tailwindcss/vite"},
            ConfigFiles: []ConfigFile{
                {Path: "tailwind.config.ts", Content: tailwindConfig},
            },
            FilePatches: []FilePatch{
                {Path: cssEntry, Insert: `@import "tailwindcss";\n`, Mode: PatchInsertAfter},
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
                "lint":        "eslint .",
                "lint:fix":    "eslint . --fix",
                "format":      "prettier --write .",
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
        "tanstack-query": {
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
        },
        "tanstack-router": {
            Name:     "TanStack Router",
            Packages: []string{"@tanstack/react-router"},
            DevPackages: []string{"@tanstack/router-plugin"},
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
            Name:     "Vuetify",
            Packages: []string{"vuetify"},
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
    case "vue", "nuxt":
        return "src/main.ts"
    case "angular":
        return "src/main.ts"
    default:
        return "src/main.tsx"
    }
}

func cssEntryFile(framework string) string {
    return "src/index.css"
}

// Config file template constants — defined below.
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestToolCatalog
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/configs.go executor/configs_test.go
git commit -m "feat: ToolSetup model and full tool catalog"
```

---

## Task 5: Step interface + SelectStep

**Files:**
- Modify: `wizard/steps/step.go`
- Modify: `wizard/steps/select.go`
- Create: `wizard/steps/select_test.go`

**Step 1: Write the failing test**

`wizard/steps/select_test.go`:
```go
package steps_test

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/wizard/steps"
)

func TestSelectStep_NavigatesDown(t *testing.T) {
    s := steps.NewSelectStep("Pick one", []string{"a", "b", "c"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
    assert.Equal(t, "b", result.Value())
}

func TestSelectStep_NavigatesUp(t *testing.T) {
    s := steps.NewSelectStep("Pick one", []string{"a", "b", "c"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyDown})
    result, _ = result.Update(tea.KeyMsg{Type: tea.KeyUp})
    assert.Equal(t, "a", result.Value())
}

func TestSelectStep_WrapsAtBounds(t *testing.T) {
    s := steps.NewSelectStep("Pick one", []string{"a", "b"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyUp})
    assert.Equal(t, "a", result.Value()) // no wrap-around above 0
}

func TestSelectStep_EnterCompletes(t *testing.T) {
    s := steps.NewSelectStep("Pick one", []string{"a", "b"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
    assert.True(t, result.IsDone())
    assert.Equal(t, "a", result.Value())
}

func TestSelectStep_ViewContainsTitle(t *testing.T) {
    s := steps.NewSelectStep("Choose framework", []string{"react", "vue"})
    assert.Contains(t, s.View(), "Choose framework")
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./wizard/steps/... -v -run TestSelectStep
```
Expected: FAIL.

**Step 3: Implement Step interface and SelectStep**

`wizard/steps/step.go`:
```go
package steps

import tea "github.com/charmbracelet/bubbletea"

type Step interface {
    Update(tea.Msg) (Step, tea.Cmd)
    View() string
    IsDone() bool
    Value() any
}
```

`wizard/steps/select.go`:
```go
package steps

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

var (
    selectedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
    cursorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
    titleStyle     = lipgloss.NewStyle().Bold(true).MarginBottom(1)
)

type SelectStep struct {
    title   string
    choices []string
    cursor  int
    done    bool
}

func NewSelectStep(title string, choices []string) Step {
    return &SelectStep{title: title, choices: choices}
}

func (s *SelectStep) Update(msg tea.Msg) (Step, tea.Cmd) {
    cp := *s
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp:
            if cp.cursor > 0 {
                cp.cursor--
            }
        case tea.KeyDown:
            if cp.cursor < len(cp.choices)-1 {
                cp.cursor++
            }
        case tea.KeyEnter:
            cp.done = true
        }
    }
    return &cp, nil
}

func (s *SelectStep) View() string {
    out := titleStyle.Render(s.title) + "\n\n"
    for i, c := range s.choices {
        cursor := "  "
        line := c
        if i == s.cursor {
            cursor = cursorStyle.Render("▶ ")
            line = selectedStyle.Render(c)
        }
        out += fmt.Sprintf("%s%s\n", cursor, line)
    }
    out += "\n" + lipgloss.NewStyle().Faint(true).Render("↑/↓ navigate • enter select")
    return out
}

func (s *SelectStep) IsDone() bool { return s.done }
func (s *SelectStep) Value() any   { return s.choices[s.cursor] }
```

**Step 4: Run test to verify it passes**

```bash
go test ./wizard/steps/... -v -run TestSelectStep
```
Expected: PASS.

**Step 5: Commit**

```bash
git add wizard/steps/
git commit -m "feat: Step interface and SelectStep with arrow navigation"
```

---

## Task 6: MultiSelectStep

**Files:**
- Modify: `wizard/steps/multiselect.go`
- Create: `wizard/steps/multiselect_test.go`

**Step 1: Write the failing test**

`wizard/steps/multiselect_test.go`:
```go
package steps_test

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/wizard/steps"
)

func TestMultiSelectStep_ToggleSelection(t *testing.T) {
    s := steps.NewMultiSelectStep("Pick tools", []string{"vitest", "playwright", "storybook"})
    // press space to select first item
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeySpace})
    vals := result.Value().([]string)
    assert.Contains(t, vals, "vitest")
}

func TestMultiSelectStep_DeselectOnSecondToggle(t *testing.T) {
    s := steps.NewMultiSelectStep("Pick tools", []string{"vitest", "playwright"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeySpace})
    result, _ = result.Update(tea.KeyMsg{Type: tea.KeySpace})
    vals := result.Value().([]string)
    assert.NotContains(t, vals, "vitest")
}

func TestMultiSelectStep_EnterCompletes(t *testing.T) {
    s := steps.NewMultiSelectStep("Pick tools", []string{"vitest"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
    assert.True(t, result.IsDone())
}

func TestMultiSelectStep_CanSelectNone(t *testing.T) {
    s := steps.NewMultiSelectStep("Pick tools", []string{"vitest"})
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
    vals := result.Value().([]string)
    assert.Empty(t, vals)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./wizard/steps/... -v -run TestMultiSelectStep
```
Expected: FAIL.

**Step 3: Implement MultiSelectStep**

`wizard/steps/multiselect.go`:
```go
package steps

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type MultiSelectStep struct {
    title    string
    choices  []string
    selected map[int]bool
    cursor   int
    done     bool
}

func NewMultiSelectStep(title string, choices []string) Step {
    return &MultiSelectStep{
        title:    title,
        choices:  choices,
        selected: make(map[int]bool),
    }
}

func (s *MultiSelectStep) Update(msg tea.Msg) (Step, tea.Cmd) {
    cp := *s
    cp.selected = make(map[int]bool, len(s.selected))
    for k, v := range s.selected {
        cp.selected[k] = v
    }
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyUp:
            if cp.cursor > 0 {
                cp.cursor--
            }
        case tea.KeyDown:
            if cp.cursor < len(cp.choices)-1 {
                cp.cursor++
            }
        case tea.KeySpace:
            cp.selected[cp.cursor] = !cp.selected[cp.cursor]
        case tea.KeyEnter:
            cp.done = true
        }
    }
    return &cp, nil
}

func (s *MultiSelectStep) View() string {
    out := titleStyle.Render(s.title) + "\n\n"
    for i, c := range s.choices {
        cursor := "  "
        if i == s.cursor {
            cursor = cursorStyle.Render("▶ ")
        }
        checkbox := "○"
        if s.selected[i] {
            checkbox = selectedStyle.Render("●")
        }
        out += fmt.Sprintf("%s%s %s\n", cursor, checkbox, c)
    }
    out += "\n" + lipgloss.NewStyle().Faint(true).Render("↑/↓ navigate • space toggle • enter confirm")
    return out
}

func (s *MultiSelectStep) IsDone() bool { return s.done }

func (s *MultiSelectStep) Value() any {
    var result []string
    for i, c := range s.choices {
        if s.selected[i] {
            result = append(result, c)
        }
    }
    if result == nil {
        return []string{}
    }
    return result
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./wizard/steps/... -v -run TestMultiSelectStep
```
Expected: PASS.

**Step 5: Commit**

```bash
git add wizard/steps/
git commit -m "feat: MultiSelectStep with space-toggle checkboxes"
```

---

## Task 7: ConfirmStep (summary screen)

**Files:**
- Modify: `wizard/steps/confirm.go`
- Create: `wizard/steps/confirm_test.go`

**Step 1: Write the failing test**

`wizard/steps/confirm_test.go`:
```go
package steps_test

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/wizard/steps"
)

func TestConfirmStep_ViewShowsAllChoices(t *testing.T) {
    cfg := config.New()
    cfg.Framework  = "react"
    cfg.Variant    = "vite"
    cfg.Linting    = "eslint-prettier"
    cfg.UILibrary  = "shadcn"
    cfg.Testing    = []string{"vitest", "playwright"}
    cfg.Tooling    = []string{"tanstack-query"}
    s := steps.NewConfirmStep(cfg)
    view := s.View()
    assert.Contains(t, view, "react")
    assert.Contains(t, view, "shadcn")
    assert.Contains(t, view, "vitest")
    assert.Contains(t, view, "tanstack-query")
}

func TestConfirmStep_EnterCompletes(t *testing.T) {
    s := steps.NewConfirmStep(config.New())
    result, _ := s.Update(tea.KeyMsg{Type: tea.KeyEnter})
    assert.True(t, result.IsDone())
    assert.Equal(t, true, result.Value())
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./wizard/steps/... -v -run TestConfirmStep
```
Expected: FAIL.

**Step 3: Implement ConfirmStep**

`wizard/steps/confirm.go`:
```go
package steps

import (
    "fmt"
    "strings"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/omaritooo/frontend-init/config"
)

var (
    labelStyle = lipgloss.NewStyle().Faint(true).Width(18)
    valueStyle = lipgloss.NewStyle().Bold(true)
)

type ConfirmStep struct {
    cfg  *config.ProjectConfig
    done bool
}

func NewConfirmStep(cfg *config.ProjectConfig) Step {
    return &ConfirmStep{cfg: cfg}
}

func (s *ConfirmStep) Update(msg tea.Msg) (Step, tea.Cmd) {
    cp := *s
    if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyEnter {
        cp.done = true
    }
    return &cp, nil
}

func (s *ConfirmStep) View() string {
    row := func(label, val string) string {
        return labelStyle.Render(label+":") + " " + valueStyle.Render(val) + "\n"
    }
    out := titleStyle.Render("Review your setup") + "\n\n"
    out += row("Mode", s.cfg.Mode)
    out += row("Package manager", s.cfg.PackageManager)
    out += row("Framework", fmt.Sprintf("%s (%s)", s.cfg.Framework, s.cfg.Variant))
    out += row("TypeScript", fmt.Sprintf("%v", s.cfg.TypeScript))
    out += row("Linting", s.cfg.Linting)
    out += row("UI library", or(s.cfg.UILibrary, "none"))
    out += row("Testing", or(strings.Join(s.cfg.Testing, ", "), "none"))
    out += row("Tooling", or(strings.Join(s.cfg.Tooling, ", "), "none"))
    out += "\n" + lipgloss.NewStyle().Faint(true).Render("enter to confirm • backspace to go back")
    return out
}

func (s *ConfirmStep) IsDone() bool { return s.done }
func (s *ConfirmStep) Value() any   { return true }

func or(a, b string) string {
    if a == "" {
        return b
    }
    return a
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./wizard/steps/... -v -run TestConfirmStep
```
Expected: PASS.

**Step 5: Commit**

```bash
git add wizard/steps/
git commit -m "feat: ConfirmStep renders summary table of all config choices"
```

---

## Task 8: ExecuteStep (progress screen)

**Files:**
- Modify: `wizard/steps/execute.go`
- Create: `wizard/steps/execute_test.go`

**Step 1: Write the failing test**

`wizard/steps/execute_test.go`:
```go
package steps_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/wizard/steps"
)

func TestExecuteStep_ViewShowsTasks(t *testing.T) {
    tasks := []steps.TaskStatus{
        {Label: "Scaffold project", State: steps.TaskDone},
        {Label: "Install packages", State: steps.TaskRunning},
        {Label: "Configure Tailwind", State: steps.TaskPending},
    }
    s := steps.NewExecuteStep(tasks)
    view := s.View()
    assert.Contains(t, view, "Scaffold project")
    assert.Contains(t, view, "Install packages")
    assert.Contains(t, view, "Configure Tailwind")
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./wizard/steps/... -v -run TestExecuteStep
```
Expected: FAIL.

**Step 3: Implement ExecuteStep**

`wizard/steps/execute.go`:
```go
package steps

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type TaskState int

const (
    TaskPending TaskState = iota
    TaskRunning
    TaskDone
    TaskFailed
)

type TaskStatus struct {
    Label string
    State TaskState
    Err   error
}

var (
    doneIcon    = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Render("✓")
    runningIcon = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("⠋")
    pendingIcon = lipgloss.NewStyle().Faint(true).Render("○")
    failedIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗")
)

type ExecuteStep struct {
    tasks []TaskStatus
    done  bool
}

func NewExecuteStep(tasks []TaskStatus) Step {
    return &ExecuteStep{tasks: tasks}
}

func (s *ExecuteStep) Update(msg tea.Msg) (Step, tea.Cmd) {
    cp := *s
    cp.tasks = make([]TaskStatus, len(s.tasks))
    copy(cp.tasks, s.tasks)
    switch msg := msg.(type) {
    case TaskProgressMsg:
        cp.tasks[msg.Index].State = msg.State
        cp.tasks[msg.Index].Err = msg.Err
        allDone := true
        for _, t := range cp.tasks {
            if t.State != TaskDone && t.State != TaskFailed {
                allDone = false
            }
        }
        cp.done = allDone
    }
    return &cp, nil
}

func (s *ExecuteStep) View() string {
    out := titleStyle.Render("Setting up your project") + "\n\n"
    for _, t := range s.tasks {
        icon := pendingIcon
        switch t.State {
        case TaskDone:    icon = doneIcon
        case TaskRunning: icon = runningIcon
        case TaskFailed:  icon = failedIcon
        }
        out += icon + "  " + t.Label + "\n"
    }
    if s.done {
        out += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true).Render("✓ All done!")
    }
    return out
}

func (s *ExecuteStep) IsDone() bool { return s.done }
func (s *ExecuteStep) Value() any   { return s.tasks }

// TaskProgressMsg is sent by the executor to update task state.
type TaskProgressMsg struct {
    Index int
    State TaskState
    Err   error
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./wizard/steps/... -v -run TestExecuteStep
```
Expected: PASS.

**Step 5: Commit**

```bash
git add wizard/steps/
git commit -m "feat: ExecuteStep renders live progress list"
```

---

## Task 9: Dynamic step builder (steps.go)

**Files:**
- Modify: `wizard/steps.go`
- Create: `wizard/steps_test.go`

**Step 1: Write the failing test**

`wizard/steps_test.go`:
```go
package wizard_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/wizard"
)

func TestBuildSteps_CustomPathHasAllScreens(t *testing.T) {
    cfg := config.New()
    cfg.Preset = "custom"
    all := wizard.BuildSteps(cfg)
    labels := stepLabels(all)
    assert.Contains(t, labels, "mode")
    assert.Contains(t, labels, "framework")
    assert.Contains(t, labels, "linting")
}

func TestBuildSteps_PresetPathSkipsToConfirm(t *testing.T) {
    cfg := config.New()
    cfg.Preset = "React Minimal"
    p := config.GetPreset("React Minimal")
    p.Apply(cfg)
    all := wizard.BuildSteps(cfg)
    labels := stepLabels(all)
    assert.NotContains(t, labels, "framework") // preset pre-fills framework
    assert.Contains(t, labels, "confirm")
}

func TestBuildSteps_AngularHasNoTanStackQuery(t *testing.T) {
    cfg := config.New()
    cfg.Framework = "angular"
    all := wizard.BuildInitialSteps(cfg) // steps after framework is resolved
    labels := stepLabels(all)
    // Tooling step choices for angular should not include tanstack-query
    // We check by building tooling options and asserting
    opts := wizard.ToolingOptions(cfg)
    for _, o := range opts {
        assert.NotEqual(t, "tanstack-router", o, "Angular should not offer TanStack Router")
    }
}

func stepLabels(steps []interface{ Label() string }) []string {
    labels := make([]string, len(steps))
    for i, s := range labels {
        _ = s
        labels[i] = steps[i].Label()
    }
    return labels
}
```

Note: `Label()` needs to be added to the `Step` interface for testing purposes. See step 3.

**Step 2: Run test to verify it fails**

```bash
go test ./wizard/... -v -run TestBuildSteps
```
Expected: FAIL.

**Step 3: Add Label() to Step interface and implement step builder**

First extend `wizard/steps/step.go`:
```go
type Step interface {
    Update(tea.Msg) (Step, tea.Cmd)
    View() string
    IsDone() bool
    Value() any
    Label() string  // identifies step for conditional logic
}
```

Add `Label()` to `SelectStep` returning `s.title` (lowercase), same for all step types.

`wizard/steps.go`:
```go
package wizard

import (
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/wizard/steps"
)

// BuildSteps constructs the ordered step slice for the wizard.
// For preset mode, only preset selection + summary steps are returned.
// For custom mode, all steps are returned.
func BuildSteps(cfg *config.ProjectConfig) []steps.Step {
    base := []steps.Step{
        steps.NewSelectStep("mode", []string{"new", "existing"}),
        steps.NewSelectStep("preset", presetNames()),
    }
    if cfg.Preset != "custom" {
        base = append(base, steps.NewConfirmStep(cfg))
        return base
    }
    return BuildInitialSteps(cfg)
}

// BuildInitialSteps returns the full custom wizard step slice.
func BuildInitialSteps(cfg *config.ProjectConfig) []steps.Step {
    s := []steps.Step{
        steps.NewSelectStep("mode", []string{"new", "existing"}),
        steps.NewSelectStep("preset", append(presetNames(), "Custom")),
        steps.NewSelectStep("package manager", []string{"npm", "pnpm", "yarn", "bun"}),
        steps.NewSelectStep("framework", []string{"react", "vue", "svelte", "angular", "astro"}),
    }
    s = append(s, variantStep(cfg)...)
    s = append(s,
        steps.NewSelectStep("typescript", []string{"yes", "no"}),
        steps.NewSelectStep("linting", []string{
            "eslint-prettier", "biome", "oxlint", "none",
        }),
        steps.NewSelectStep("ui library", UILibraryOptions(cfg)),
        steps.NewMultiSelectStep("testing", TestingOptions(cfg)),
        steps.NewMultiSelectStep("tooling", ToolingOptions(cfg)),
        steps.NewConfirmStep(cfg),
    )
    return s
}

func variantStep(cfg *config.ProjectConfig) []steps.Step {
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

// UILibraryOptions returns the UI library choices filtered by framework.
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

// ToolingOptions returns tooling choices filtered by framework.
func ToolingOptions(cfg *config.ProjectConfig) []string {
    switch cfg.Framework {
    case "react":
        base := []string{"tanstack-query", "zustand", "jotai", "redux-toolkit", "rhf-zod", "zod", "axios", "i18next"}
        if cfg.Variant != "nextjs" {
            base = append([]string{"tanstack-router", "react-router-v7"}, base...)
        } else {
            base = append(base, "trpc")
        }
        return base
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

func presetNames() []string {
    presets := config.AllPresets()
    names := make([]string, len(presets))
    for i, p := range presets {
        names[i] = p.Name
    }
    return names
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./wizard/... -v -run TestBuildSteps
```
Expected: PASS.

**Step 5: Commit**

```bash
git add wizard/
git commit -m "feat: dynamic step builder with framework-filtered options"
```

---

## Task 10: Root wizard model

**Files:**
- Modify: `wizard/model.go`
- Create: `wizard/model_test.go`

**Step 1: Write the failing test**

`wizard/model_test.go`:
```go
package wizard_test

import (
    "testing"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/wizard"
)

func TestModel_AdvancesOnStepComplete(t *testing.T) {
    cfg := config.New()
    m := wizard.New(cfg)
    assert.Equal(t, 0, m.Cursor())
    // press enter on first step
    newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
    wm := newM.(wizard.Model)
    assert.Equal(t, 1, wm.Cursor())
}

func TestModel_BacktrackOnEsc(t *testing.T) {
    cfg := config.New()
    m := wizard.New(cfg)
    newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // advance to step 1
    wm := newM.(wizard.Model)
    newM, _ = wm.Update(tea.KeyMsg{Type: tea.KeyEsc})   // go back
    wm = newM.(wizard.Model)
    assert.Equal(t, 0, wm.Cursor())
}

func TestModel_DoesNotGoBeforeFirst(t *testing.T) {
    cfg := config.New()
    m := wizard.New(cfg)
    newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
    wm := newM.(wizard.Model)
    assert.Equal(t, 0, wm.Cursor())
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./wizard/... -v -run TestModel
```
Expected: FAIL.

**Step 3: Implement root model**

`wizard/model.go`:
```go
package wizard

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/wizard/steps"
)

type Model struct {
    stepList []steps.Step
    cursor   int
    cfg      *config.ProjectConfig
}

func New(cfg *config.ProjectConfig) Model {
    return Model{
        stepList: BuildInitialSteps(cfg),
        cfg:      cfg,
    }
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        switch key.Type {
        case tea.KeyCtrlC:
            return m, tea.Quit
        case tea.KeyEsc:
            if m.cursor > 0 {
                m.cursor--
            }
            return m, nil
        }
    }

    current := m.stepList[m.cursor]
    newStep, cmd := current.Update(msg)
    m.stepList[m.cursor] = newStep

    if newStep.IsDone() {
        m.applyStepValue(newStep)
        if m.cursor < len(m.stepList)-1 {
            m.cursor++
            // rebuild conditional steps when framework is set
            if newStep.Label() == "framework" {
                m.stepList = rebuildAfterFramework(m.stepList, m.cursor, m.cfg)
            }
        } else {
            return m, tea.Quit
        }
    }
    return m, cmd
}

func (m Model) View() string {
    return m.stepList[m.cursor].View()
}

func (m Model) Cursor() int { return m.cursor }

// applyStepValue writes the step's result into ProjectConfig.
func (m *Model) applyStepValue(s steps.Step) {
    val := s.Value()
    switch s.Label() {
    case "mode":
        m.cfg.Mode = val.(string)
    case "preset":
        name := val.(string)
        if name != "Custom" {
            if p := config.GetPreset(name); p != nil {
                p.Apply(m.cfg)
            }
        }
        m.cfg.Preset = name
    case "package manager":
        m.cfg.PackageManager = val.(string)
    case "framework":
        m.cfg.Framework = val.(string)
    case "variant":
        m.cfg.Variant = val.(string)
    case "typescript":
        m.cfg.TypeScript = val.(string) == "yes"
    case "linting":
        m.cfg.Linting = val.(string)
    case "ui library":
        m.cfg.UILibrary = val.(string)
    case "testing":
        m.cfg.Testing = val.([]string)
    case "tooling":
        m.cfg.Tooling = val.([]string)
    }
}

// rebuildAfterFramework replaces the tail of the step list with
// framework-appropriate options from the current cursor onward.
func rebuildAfterFramework(current []steps.Step, from int, cfg *config.ProjectConfig) []steps.Step {
    head := make([]steps.Step, from)
    copy(head, current)
    tail := []steps.Step{
        NewSelectStep("typescript", []string{"yes", "no"}),
        NewSelectStep("linting", []string{"eslint-prettier", "biome", "oxlint", "none"}),
        NewSelectStep("ui library", UILibraryOptions(cfg)),
        NewMultiSelectStep("testing", TestingOptions(cfg)),
        NewMultiSelectStep("tooling", ToolingOptions(cfg)),
        NewConfirmStep(cfg),
    }
    return append(head, tail...)
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./wizard/... -v -run TestModel
```
Expected: PASS.

**Step 5: Commit**

```bash
git add wizard/
git commit -m "feat: root Bubbletea wizard model with step navigation and config wiring"
```

---

## Task 11: Package manager detection

**Files:**
- Create: `executor/detect.go`
- Create: `executor/detect_test.go`

**Step 1: Write the failing test**

`executor/detect_test.go`:
```go
package executor_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/executor"
)

func TestDetectPackageManager_FromLockfile(t *testing.T) {
    tests := []struct{
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
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestDetectPackageManager
```
Expected: FAIL.

**Step 3: Implement detection**

`executor/detect.go`:
```go
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestDetectPackageManager
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/detect.go executor/detect_test.go
git commit -m "feat: package manager auto-detection from lockfiles"
```

---

## Task 12: CommandRunner interface + scaffold commands

**Files:**
- Modify: `executor/commands.go`
- Create: `executor/commands_test.go`

**Step 1: Write the failing test**

`executor/commands_test.go`:
```go
package executor_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/executor"
)

func TestScaffoldCmd_ReactVite(t *testing.T) {
    cfg := &config.ProjectConfig{
        Framework: "react", Variant: "vite",
        PackageManager: "npm", TypeScript: true,
    }
    cmd := executor.ScaffoldCommand(cfg, "my-app")
    assert.Equal(t, "npm", cmd[0])
    assert.Contains(t, cmd, "create")
    assert.Contains(t, cmd, "vite@latest")
}

func TestScaffoldCmd_NextJS(t *testing.T) {
    cfg := &config.ProjectConfig{
        Framework: "react", Variant: "nextjs",
        PackageManager: "pnpm",
    }
    cmd := executor.ScaffoldCommand(cfg, "my-app")
    assert.Equal(t, "pnpm", cmd[0])
    assert.Contains(t, cmd, "dlx")
    assert.Contains(t, cmd, "create-next-app@latest")
}

func TestScaffoldCmd_Angular(t *testing.T) {
    cfg := &config.ProjectConfig{
        Framework: "angular", Variant: "angular-cli",
        PackageManager: "npm",
    }
    cmd := executor.ScaffoldCommand(cfg, "my-app")
    assert.Contains(t, cmd, "ng")
    assert.Contains(t, cmd, "new")
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestScaffoldCmd
```
Expected: FAIL.

**Step 3: Implement CommandRunner and scaffold commands**

`executor/commands.go`:
```go
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
    ts := ""
    if cfg.TypeScript {
        ts = "--typescript"
    }
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
        return []string{"npx", "@angular/cli@latest", "new", projectName, ts}
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestScaffoldCmd
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/commands.go executor/commands_test.go
git commit -m "feat: scaffold and install commands for all frameworks"
```

---

## Task 13: Config file writer

**Files:**
- Modify: `executor/configs.go` (add WriteConfigFiles)
- Create: `executor/configs_write_test.go`

**Step 1: Write the failing test**

`executor/configs_write_test.go`:
```go
package executor_test

import (
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/executor"
)

func TestWriteConfigFiles_CreatesFile(t *testing.T) {
    dir := t.TempDir()
    files := []executor.ConfigFile{
        {Path: "eslint.config.js", Content: "export default []"},
        {Path: ".prettierrc", Content: `{"semi": false}`},
    }
    err := executor.WriteConfigFiles(dir, files)
    assert.NoError(t, err)
    for _, f := range files {
        data, err := os.ReadFile(filepath.Join(dir, f.Path))
        assert.NoError(t, err)
        assert.Equal(t, f.Content, string(data))
    }
}

func TestWriteConfigFiles_CreatesSubdirs(t *testing.T) {
    dir := t.TempDir()
    files := []executor.ConfigFile{
        {Path: "src/config/env.ts", Content: "export const env = {}"},
    }
    err := executor.WriteConfigFiles(dir, files)
    assert.NoError(t, err)
    _, err = os.Stat(filepath.Join(dir, "src/config/env.ts"))
    assert.NoError(t, err)
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestWriteConfigFiles
```
Expected: FAIL.

**Step 3: Implement WriteConfigFiles**

Add to `executor/configs.go`:
```go
import (
    "os"
    "path/filepath"
)

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
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestWriteConfigFiles
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/configs.go executor/configs_write_test.go
git commit -m "feat: WriteConfigFiles writes tool configs with auto dir creation"
```

---

## Task 14: File patcher

**Files:**
- Modify: `executor/patches.go`
- Create: `executor/patches_test.go`

**Step 1: Write the failing test**

`executor/patches_test.go`:
```go
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
    assert.NoError(t, err) // gracefully skip missing files
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestPatchFile
```
Expected: FAIL.

**Step 3: Implement file patcher**

`executor/patches.go`:
```go
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
        return nil // skip gracefully
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestPatchFile
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/patches.go executor/patches_test.go
git commit -m "feat: file patcher for post-install modifications to existing files"
```

---

## Task 15: Package.json script merger

**Files:**
- Create: `executor/pkgjson.go`
- Create: `executor/pkgjson_test.go`

**Step 1: Write the failing test**

`executor/pkgjson_test.go`:
```go
package executor_test

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/executor"
)

func TestMergeScripts_AddsNewScripts(t *testing.T) {
    dir := t.TempDir()
    pkg := map[string]any{
        "name": "my-app",
        "scripts": map[string]any{"dev": "vite"},
    }
    data, _ := json.Marshal(pkg)
    _ = os.WriteFile(filepath.Join(dir, "package.json"), data, 0644)

    err := executor.MergeScripts(dir, map[string]string{
        "lint":   "eslint .",
        "format": "prettier --write .",
    })
    assert.NoError(t, err)

    result, _ := os.ReadFile(filepath.Join(dir, "package.json"))
    var out map[string]any
    _ = json.Unmarshal(result, &out)
    scripts := out["scripts"].(map[string]any)
    assert.Equal(t, "eslint .", scripts["lint"])
    assert.Equal(t, "vite", scripts["dev"]) // existing preserved
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestMergeScripts
```
Expected: FAIL.

**Step 3: Implement MergeScripts**

`executor/pkgjson.go`:
```go
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
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestMergeScripts
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/pkgjson.go executor/pkgjson_test.go
git commit -m "feat: MergeScripts merges tool scripts into package.json"
```

---

## Task 16: Executor orchestration

**Files:**
- Modify: `executor/executor.go`
- Create: `executor/executor_test.go`

**Step 1: Write the failing test**

`executor/executor_test.go`:
```go
package executor_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/executor"
)

type mockRunner struct {
    ran [][]string
}

func (m *mockRunner) Run(dir, name string, args ...string) error {
    m.ran = append(m.ran, append([]string{name}, args...))
    return nil
}

func TestExecutor_RunsScaffoldForNewProject(t *testing.T) {
    cfg := config.New()
    cfg.Mode      = "new"
    cfg.Framework = "react"
    cfg.Variant   = "vite"
    cfg.Linting   = "none"
    cfg.TypeScript = true

    runner := &mockRunner{}
    dir    := t.TempDir()
    ex     := executor.New(cfg, runner, dir)
    tasks  := ex.Tasks()

    assert.True(t, len(tasks) > 0)
    assert.Equal(t, "Scaffold project", tasks[0].Label)
}

func TestExecutor_SkipsScaffoldForExistingProject(t *testing.T) {
    cfg := config.New()
    cfg.Mode      = "existing"
    cfg.Framework = "react"
    cfg.Variant   = "vite"
    cfg.Linting   = "none"

    runner := &mockRunner{}
    dir    := t.TempDir()
    ex     := executor.New(cfg, runner, dir)
    tasks  := ex.Tasks()

    for _, t2 := range tasks {
        assert.NotEqual(t, "Scaffold project", t2.Label)
    }
}
```

**Step 2: Run test to verify it fails**

```bash
go test ./executor/... -v -run TestExecutor
```
Expected: FAIL.

**Step 3: Implement Executor**

`executor/executor.go`:
```go
package executor

import (
    "os/exec"
    "strings"
    "github.com/omaritooo/frontend-init/config"
    wsteps "github.com/omaritooo/frontend-init/wizard/steps"
)

// Task is a labelled unit of work with an execution function.
type Task struct {
    Label string
    Run   func() error
}

// Executor builds and runs the setup task list.
type Executor struct {
    cfg        *config.ProjectConfig
    runner     CommandRunner
    projectDir string
    projectName string
}

func New(cfg *config.ProjectConfig, runner CommandRunner, projectDir string) *Executor {
    return &Executor{cfg: cfg, runner: runner, projectDir: projectDir}
}

func (e *Executor) SetProjectName(name string) { e.projectName = name }

// Tasks returns the ordered list of setup tasks derived from config.
func (e *Executor) Tasks() []Task {
    var tasks []Task

    // 1. Scaffold (new projects only)
    if e.cfg.IsNewProject() {
        cmd := ScaffoldCommand(e.cfg, e.projectName)
        if cmd != nil {
            tasks = append(tasks, Task{
                Label: "Scaffold project",
                Run: func() error {
                    return e.runner.Run(".", cmd[0], cmd[1:]...)
                },
            })
        }
    }

    // Collect all selected tools
    tools := e.selectedTools()

    // 2. Install packages (batched)
    var pkgs, devPkgs []string
    for _, t := range tools {
        pkgs = append(pkgs, t.Packages...)
        devPkgs = append(devPkgs, t.DevPackages...)
    }
    if len(pkgs) > 0 {
        p := pkgs
        tasks = append(tasks, Task{
            Label: "Install dependencies",
            Run: func() error {
                cmd := InstallCmd(e.cfg.PackageManager, false, p)
                return e.runner.Run(e.projectDir, cmd[0], cmd[1:]...)
            },
        })
    }
    if len(devPkgs) > 0 {
        d := devPkgs
        tasks = append(tasks, Task{
            Label: "Install dev dependencies",
            Run: func() error {
                cmd := InstallCmd(e.cfg.PackageManager, true, d)
                return e.runner.Run(e.projectDir, cmd[0], cmd[1:]...)
            },
        })
    }

    // 3. Write config files
    for _, t := range tools {
        if len(t.ConfigFiles) > 0 {
            tool := t
            tasks = append(tasks, Task{
                Label: "Configure " + tool.Name,
                Run: func() error {
                    return WriteConfigFiles(e.projectDir, tool.ConfigFiles)
                },
            })
        }
    }

    // 4. Patch existing files
    for _, t := range tools {
        if len(t.FilePatches) > 0 {
            tool := t
            tasks = append(tasks, Task{
                Label: "Patch files for " + tool.Name,
                Run: func() error {
                    return ApplyPatches(e.projectDir, tool.FilePatches)
                },
            })
        }
    }

    // 5. Merge package.json scripts
    allScripts := make(map[string]string)
    for _, t := range tools {
        for k, v := range t.Scripts {
            allScripts[k] = v
        }
    }
    if len(allScripts) > 0 {
        s := allScripts
        tasks = append(tasks, Task{
            Label: "Update package.json scripts",
            Run: func() error {
                return MergeScripts(e.projectDir, s)
            },
        })
    }

    // 6. Post-install commands
    for _, t := range tools {
        for _, postCmd := range t.PostInstallCmds {
            pc := postCmd
            tool := t
            args := strings.Fields(pc)
            tasks = append(tasks, Task{
                Label: tool.Name + ": " + pc,
                Run: func() error {
                    return e.runner.Run(e.projectDir, args[0], args[1:]...)
                },
            })
        }
    }

    return tasks
}

// selectedTools returns ToolSetup instances for all selected tools in config.
func (e *Executor) selectedTools() []ToolSetup {
    var tools []ToolSetup
    add := func(key string) {
        if s := GetToolSetup(key, e.cfg.Framework); s != nil {
            tools = append(tools, *s)
        }
    }
    if e.cfg.Linting != "none" && e.cfg.Linting != "" {
        add(e.cfg.Linting)
    }
    if e.cfg.UILibrary != "none" && e.cfg.UILibrary != "" {
        if e.cfg.UILibrary != "tailwind-only" {
            add("tailwind") // UI libs that need tailwind get it first
        }
        add(e.cfg.UILibrary)
    }
    for _, t := range e.cfg.Testing {
        add(t)
    }
    for _, t := range e.cfg.Tooling {
        add(t)
    }
    return tools
}

// RealRunner implements CommandRunner using os/exec.
type RealRunner struct{}

func (r *RealRunner) Run(dir, name string, args ...string) error {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir
    cmd.Stdout = nil
    cmd.Stderr = nil
    return cmd.Run()
}

// ToWizardTasks converts executor Tasks to wizard ExecuteStep TaskStatus slice.
func ToWizardTasks(tasks []Task) []wsteps.TaskStatus {
    result := make([]wsteps.TaskStatus, len(tasks))
    for i, t := range tasks {
        result[i] = wsteps.TaskStatus{Label: t.Label, State: wsteps.TaskPending}
    }
    return result
}
```

**Step 4: Run test to verify it passes**

```bash
go test ./executor/... -v -run TestExecutor
```
Expected: PASS.

**Step 5: Commit**

```bash
git add executor/executor.go executor/executor_test.go
git commit -m "feat: executor orchestration — scaffold, install, config, patch, post-install"
```

---

## Task 17: Wire Cobra commands

**Files:**
- Modify: `cmd/root.go`
- Modify: `cmd/init.go`

**Step 1: Update root.go**

`cmd/root.go` — strip template boilerplate, set meaningful metadata:
```go
package cmd

import (
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "frontend-init",
    Short: "Scaffold and configure frontend projects",
    Long:  "An interactive TUI wizard to scaffold React, Vue, Svelte, Angular, and Astro projects with your preferred tooling.",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    rootCmd.AddCommand(initCmd)
}
```

**Step 2: Implement init.go**

`cmd/init.go`:
```go
package cmd

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/spf13/cobra"
    "github.com/omaritooo/frontend-init/config"
    "github.com/omaritooo/frontend-init/executor"
    "github.com/omaritooo/frontend-init/wizard"
    wsteps "github.com/omaritooo/frontend-init/wizard/steps"
)

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Start the interactive setup wizard",
    RunE:  runInit,
}

func runInit(_ *cobra.Command, _ []string) error {
    cfg := config.New()

    // Auto-detect package manager for existing projects
    if wd, err := os.Getwd(); err == nil {
        cfg.PackageManager = executor.DetectPackageManager(wd)
    }

    m := wizard.New(cfg)
    p := tea.NewProgram(m, tea.WithAltScreen())
    finalModel, err := p.Run()
    if err != nil {
        return err
    }

    wm, ok := finalModel.(wizard.Model)
    if !ok {
        return fmt.Errorf("unexpected model type")
    }

    finalCfg := wm.Config()
    wd, _ := os.Getwd()
    projectDir := wd

    runner := &executor.RealRunner{}
    ex := executor.New(finalCfg, runner, projectDir)
    tasks := ex.Tasks()
    wizardTasks := executor.ToWizardTasks(tasks)

    // Run execution phase as a second Bubbletea program
    execStep := wsteps.NewExecuteStep(wizardTasks)
    execModel := wizard.NewExecuteModel(execStep, tasks)
    ep := tea.NewProgram(execModel, tea.WithAltScreen())
    _, err = ep.Run()
    return err
}
```

**Step 3: Add Config() accessor to wizard.Model**

In `wizard/model.go`, add:
```go
func (m Model) Config() *config.ProjectConfig { return m.cfg }
```

**Step 4: Create NewExecuteModel in wizard/model.go**

```go
type ExecuteModel struct {
    step  steps.Step
    tasks []executor.Task
    index int
}

func NewExecuteModel(step steps.Step, tasks []executor.Task) ExecuteModel {
    return ExecuteModel{step: step, tasks: tasks}
}

func (e ExecuteModel) Init() tea.Cmd {
    return runNextTask(e.tasks, 0)
}

func (e ExecuteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case steps.TaskProgressMsg:
        newStep, _ := e.step.Update(msg)
        e.step = newStep
        if msg.State == steps.TaskDone && e.index+1 < len(e.tasks) {
            e.index++
            return e, runNextTask(e.tasks, e.index)
        }
        if e.step.IsDone() {
            return e, tea.Quit
        }
    case tea.KeyMsg:
        if msg.Type == tea.KeyCtrlC {
            return e, tea.Quit
        }
    }
    return e, nil
}

func (e ExecuteModel) View() string { return e.step.View() }

func runNextTask(tasks []executor.Task, idx int) tea.Cmd {
    return func() tea.Msg {
        err := tasks[idx].Run()
        state := steps.TaskDone
        if err != nil {
            state = steps.TaskFailed
        }
        return steps.TaskProgressMsg{Index: idx, State: state, Err: err}
    }
}
```

Note: this introduces an import cycle between `wizard` and `executor`. Resolve by moving `executor.Task` into a shared `types` package, or pass tasks as `[]func() error` with labels. The implementation plan allows for this refactor during implementation.

**Step 5: Build and verify it compiles**

```bash
go build ./...
```
Expected: no errors.

**Step 6: Commit**

```bash
git add cmd/
git commit -m "feat: wire Cobra commands to wizard and executor"
```

---

## Task 18: Smoke test the binary

**Step 1: Build binary**

```bash
go build -o frontend-init .
```

**Step 2: Run help**

```bash
./frontend-init --help
./frontend-init init --help
```
Expected: help text with "Scaffold and configure frontend projects".

**Step 3: Run full suite**

```bash
go test ./... -v
```
Expected: all tests pass.

**Step 4: Final commit**

```bash
git add .
git commit -m "chore: all tests passing, binary builds cleanly"
```

---

## Notes

- The import cycle between `wizard` and `executor` (Task 17) may require introducing a `types` package with `Task` defined there, imported by both. Address this when the cycle surfaces during `go build`.
- `survey/v2` is in `go.mod` but is not used — run `go mod tidy` after all tasks to clean it up.
- Angular's `ng new` and `ng add` commands require `@angular/cli` to be globally installed; document this as a prerequisite in the README.
- All Bubbletea step tests use `tea.KeyMsg{Type: tea.KeyEnter}` etc — verify these key constants match the installed bubbletea version.
