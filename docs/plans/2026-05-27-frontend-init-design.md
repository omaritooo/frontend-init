# frontend-init Design

**Date:** 2026-05-27
**Status:** Approved

## Overview

A Go CLI tool that scaffolds or configures frontend projects through an interactive Bubbletea TUI wizard. Supports all major frontend frameworks with framework-filtered tooling menus and named presets for common stacks.

---

## Commands

```
frontend-init           # runs wizard (alias for init)
frontend-init init      # explicit subcommand
```

---

## Wizard Flow (Linear)

Screens advance with `Enter`/`Tab`, retreat with `Backspace`/`Esc`. Framework-conditional screens are inserted or skipped dynamically based on prior answers.

1. **Mode** — New project | Configure existing directory
2. **Preset or Custom** — Pick a named preset (pre-fills all choices) or go step-by-step. Preset users land on Summary and can still override individual choices.
3. **Package manager** — npm | pnpm | yarn | bun (auto-detected from lockfiles in existing projects)
4. **Framework** — React | Vue | Svelte | Angular | Astro
5. **Variant** *(conditional on framework)*
   - React → Vite SPA | Next.js
   - Vue → Vite SPA | Nuxt
   - Svelte → Vite | SvelteKit
   - Angular → Angular CLI | Analog (AnalogJS)
   - Astro → Static | SSR
6. **TypeScript** — Yes | No
7. **Linting / Formatting** — None | ESLint+Prettier | Biome | Oxlint+Prettier
8. **UI Library** *(framework-filtered — see table below)*
9. **Testing** *(framework-filtered — see table below)*
10. **Tooling** *(framework-filtered — multi-select — see table below)*
11. **Summary / Confirm** — shows all choices as a table; Back to revise, Enter to execute
12. **Execution** — animated progress list showing each setup task as it runs

---

## Framework-Filtered Options

### UI Libraries

| Framework | Options |
|-----------|---------|
| React / Next.js | None · Tailwind only · shadcn/ui · MUI · Mantine · Chakra UI · Ant Design · PrimeReact · DaisyUI |
| Vue / Nuxt | None · Tailwind only · Vuetify · PrimeVue · Naive UI · DaisyUI |
| Angular / Analog | None · Tailwind only · Angular Material · PrimeNG · NG-Zorro |
| Svelte / SvelteKit | None · Tailwind only · shadcn-svelte · Skeleton UI · DaisyUI |
| Astro | None · Tailwind only · DaisyUI · shadcn/ui *(if React islands enabled)* |

### Testing

| Framework | Options |
|-----------|---------|
| React / Vue / Svelte / Astro | None · Vitest · Jest · Testing Library · Playwright · Cypress · Storybook |
| Angular / Analog | None · Jest · Playwright · Cypress |

### Tooling

| Framework | Options |
|-----------|---------|
| React / Next.js | TanStack Query · TanStack Router *(not Next.js)* · React Router v7 *(not Next.js)* · Zustand · Jotai · Redux Toolkit · React Hook Form + Zod · tRPC *(Next.js only)* · Axios/ky · i18next |
| Vue / Nuxt | TanStack Query · Pinia · VeeValidate + Zod · Axios/ky · vue-i18n |
| Angular / Analog | NgRx Signal Store · Axios/ky |
| Svelte / SvelteKit | TanStack Query · Superforms + Zod |
| Astro | Nanostores · Zod |

---

## Presets

| Name | Framework | Linting | UI | Testing | Tooling |
|------|-----------|---------|-----|---------|---------|
| React Minimal | Vite+React+TS | ESLint+Prettier | — | Vitest | — |
| React Full SPA | Vite+React+TS | ESLint+Prettier | shadcn/ui | Vitest, Testing Library, Playwright | TanStack Query, TanStack Router, Zustand, RHF+Zod |
| Next.js Standard | Next.js+TS | ESLint+Prettier | shadcn/ui | Vitest | TanStack Query, Zod |
| T3 Stack | Next.js+TS | ESLint+Prettier | shadcn/ui | Vitest | tRPC, TanStack Query, Zod |
| Vue Minimal | Vite+Vue+TS | ESLint+Prettier | — | Vitest | — |
| Vue Full SPA | Vite+Vue+TS | ESLint+Prettier | DaisyUI | Vitest | Pinia, TanStack Query, VeeValidate+Zod |
| Nuxt Standard | Nuxt+TS | ESLint+Prettier | — | Vitest | — |
| Angular Minimal | Angular CLI+TS | ESLint+Prettier | — | Jest | — |
| Angular Enterprise | Angular CLI+TS | ESLint+Prettier | Angular Material | Jest, Playwright | NgRx Signal Store |
| Astro Content Site | Astro+TS | ESLint+Prettier | Tailwind+DaisyUI | Vitest | — |
| Astro Islands | Astro+React+TS | ESLint+Prettier | shadcn/ui | Vitest | TanStack Query |

---

## Technical Architecture

### File Structure

```
cmd/
  root.go
  init.go            # entry point — launches Bubbletea app
wizard/
  model.go           # root Bubbletea model, step navigation
  steps/
    select.go        # single-choice step
    multiselect.go   # checkbox multi-select step
    input.go         # text input step
    confirm.go       # summary/review screen
    execute.go       # execution progress screen
  steps.go           # builds the ordered step slice from config
config/
  config.go          # ProjectConfig struct
  presets.go         # preset definitions
executor/
  executor.go        # orchestrates the setup sequence
  commands.go        # scaffold/install commands per framework+variant
  configs.go         # config file templates per tool
  patches.go         # existing-file patching logic (main.tsx, index.css, etc.)
```

### Bubbletea Step Interface

Each wizard screen implements:

```go
type Step interface {
    Update(tea.Msg) (Step, tea.Cmd)
    View() string
    IsDone() bool
    Value() any   // the user's selection(s) for this step
}
```

Step types: `SelectStep`, `MultiSelectStep`, `InputStep`, `ConfirmStep`, `ExecuteStep`.

The root model holds a `[]Step` slice and a cursor index. Framework-conditional steps are inserted/removed when the framework choice resolves.

### Config Struct

```go
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
```

### Tool Setup Model

Every tool is defined as a `ToolSetup` — not just a package name but a full setup sequence:

```go
type ToolSetup struct {
    Name            string
    Packages        []string          // runtime npm packages
    DevPackages     []string          // npm -D packages
    ConfigFiles     []ConfigFile      // files to create from templates
    FilePatches     []FilePatch       // modifications to existing files
    PostInstallCmds []string          // extra CLI commands after install
    Scripts         map[string]string // package.json script entries to add
}
```

Examples:

| Tool | Packages | Config files | File patches | Post-install |
|------|----------|--------------|--------------|--------------|
| Tailwind | `tailwindcss @tailwindcss/vite` | `tailwind.config.ts` | Add `@import "tailwindcss"` to `index.css` | — |
| shadcn/ui | `shadcn` | — | — | `npx shadcn init` |
| ESLint+Prettier | `eslint prettier eslint-config-prettier` | `eslint.config.js` `.prettierrc` `.prettierignore` | Add lint/format scripts | — |
| Biome | `@biomejs/biome` | — | Add lint/format scripts | `npx biome init` |
| Vitest | `vitest @vitest/ui` | `vitest.config.ts` | Add `test` script | — |
| Playwright | `@playwright/test` | `playwright.config.ts` | Add `e2e` script | `npx playwright install` |
| TanStack Query | `@tanstack/react-query` | — | Wrap `main.tsx` with `<QueryClientProvider>` | — |
| Pinia | `pinia` | — | Patch `main.ts`: `app.use(createPinia())` | — |
| NgRx Signal Store | `@ngrx/signals` | — | — | — |

### Execution Progress UI

```
✓  Scaffold project     (npm create vite@latest my-app)
✓  Install packages
⠋  Configure Tailwind CSS
   Configure ESLint + Prettier
   Patch main.tsx
   Run shadcn init
   Run playwright install
```

Each task line updates in-place via Bubbletea. On completion, a success summary prints the project path and suggested next steps.

---

## Key Design Decisions

- **Step slice is dynamic**: steps for framework variant, UI library, testing, and tooling are only added after the framework is chosen, keeping the wizard contextually relevant.
- **Presets pre-fill but don't lock**: selecting a preset populates `ProjectConfig` defaults; the user still sees the Summary screen and can go back to override any individual choice.
- **Multi-step tool setup**: every tool carries its own full setup sequence (install → config files → file patches → post-install commands) so the executor never needs to hard-code per-tool logic outside of `configs.go`.
- **Existing project mode**: skips the scaffold step, auto-detects package manager from lockfiles, and detects the framework from `package.json` dependencies to pre-select sensible defaults.
