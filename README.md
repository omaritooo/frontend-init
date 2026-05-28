# frontend-init

An interactive terminal wizard that scaffolds or configures frontend projects with your preferred framework, tooling, and libraries — fully automated from a single command.

---

## Features

- **Interactive TUI wizard** — arrow-key navigation, animated progress, summary screen before execution
- **New project or existing directory** — scaffold from scratch or add tooling to an existing project
- **5 frameworks** — React, Vue, Svelte, Angular, Astro
- **11 named presets** — one Enter to pre-fill an entire stack (still editable)
- **Framework-filtered options** — UI libraries, testing tools, and extra tooling are scoped to what actually works with your chosen framework
- **Fully automated setup** — installs packages, writes config files, patches `vite.config.ts` / `tsconfig.json` / entry files, and runs post-install commands (e.g. `npx shadcn@latest init`)
- **Animated execution progress** — live spinner per task, ticks while work runs

---

## Installation

```bash
go install github.com/omaritooo/frontend-init@latest
```

Or build from source:

```bash
git clone https://github.com/omaritooo/frontend-init
cd frontend-init
go build -o frontend-init .
```

---

## Usage

```bash
frontend-init        # runs the wizard
frontend-init init   # explicit subcommand (same thing)
```

### Navigation

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move cursor |
| `Space` | Toggle item (multi-select steps) |
| `Enter` | Confirm selection / advance |
| `Esc` | Go back one step |
| `Ctrl+C` | Quit |

---

## Wizard Flow

```
1. Mode            — New project  │  Configure existing directory
2. Project name    — (new projects only)
3. Preset          — Pick a named preset or go step-by-step
4. Package manager — npm │ pnpm │ yarn │ bun
5. Framework       — React │ Vue │ Svelte │ Angular │ Astro
6. Variant         — e.g. Vite SPA │ Next.js (framework-specific)
7. TypeScript      — Yes │ No
8. Linting         — ESLint+Prettier │ Biome │ Oxlint+Prettier │ None
9. UI library      — (framework-filtered, see below)
10. shadcn theme    — (only shown when shadcn/ui or shadcn-svelte is selected)
                      zinc │ slate │ gray │ neutral │ stone │ red │ rose │ orange │ green │ blue │ violet │ yellow
11. Testing         — (framework-filtered, multi-select)
12. Tooling         — (framework-filtered, multi-select)
13. Summary         — Review all choices, press Enter to execute
14. Execution       — Live progress with animated spinner per task
```

---

## Presets

Select a preset at step 3 to pre-fill the entire stack. You can still go back and override any individual choice.

| Preset | Framework | UI | Testing | Tooling |
|--------|-----------|-----|---------|---------|
| **React Minimal** | Vite + React + TS | — | Vitest | — |
| **React Full SPA** | Vite + React + TS | shadcn/ui | Vitest, Testing Library, Playwright | TanStack Query, TanStack Router, Zustand, RHF+Zod |
| **Next.js Standard** | Next.js + TS | shadcn/ui | Vitest | TanStack Query, Zod |
| **T3 Stack** | Next.js + TS | shadcn/ui | Vitest | tRPC, TanStack Query, Zod |
| **Vue Minimal** | Vite + Vue + TS | — | Vitest | — |
| **Vue Full SPA** | Vite + Vue + TS | DaisyUI | Vitest | Pinia, TanStack Query, VeeValidate+Zod |
| **Nuxt Standard** | Nuxt + TS | — | Vitest | — |
| **Angular Minimal** | Angular CLI + TS | — | Jest | — |
| **Angular Enterprise** | Angular CLI + TS | Angular Material | Jest, Playwright | NgRx Signal Store |
| **Astro Content Site** | Astro + TS | Tailwind + DaisyUI | Vitest | — |
| **Astro Islands** | Astro + React + TS | shadcn/ui | Vitest | TanStack Query |

---

## Framework Options

### Variants

| Framework | Variants |
|-----------|----------|
| React | Vite SPA, Next.js |
| Vue | Vite SPA, Nuxt |
| Svelte | Vite, SvelteKit |
| Angular | Angular CLI, Analog (AnalogJS) |
| Astro | Static, SSR |

### UI Libraries

| Framework | Options |
|-----------|---------|
| React / Next.js | shadcn/ui, MUI, Mantine, Chakra UI, Ant Design, PrimeReact, DaisyUI, Tailwind only |
| Vue / Nuxt | Vuetify, PrimeVue, Naive UI, DaisyUI, Tailwind only |
| Angular / Analog | Angular Material, PrimeNG, NG-Zorro, Tailwind only |
| Svelte / SvelteKit | shadcn-svelte, Skeleton UI, DaisyUI, Tailwind only |
| Astro | shadcn/ui, DaisyUI, Tailwind only |

### Testing

| Framework | Options |
|-----------|---------|
| React, Vue, Svelte, Astro | Vitest, Jest, Testing Library, Playwright, Cypress, Storybook |
| Angular / Analog | Jest, Playwright, Cypress |

### Tooling (multi-select)

| Framework | Options |
|-----------|---------|
| React (Vite) | TanStack Query, TanStack Router, React Router v7, Zustand, Jotai, Redux Toolkit, RHF+Zod, Zod, Axios, i18next |
| React (Next.js) | TanStack Query, Zustand, Jotai, Redux Toolkit, RHF+Zod, Zod, Axios, i18next, tRPC |
| Vue / Nuxt | Pinia, TanStack Query, VeeValidate+Zod, Axios, vue-i18n |
| Angular / Analog | NgRx Signal Store, Axios |
| Svelte / SvelteKit | TanStack Query, Superforms+Zod |
| Astro | Nanostores, Zod |

---

## What Gets Automated

Every tool carries its complete setup sequence — not just `npm install`:

### Tailwind CSS (v4)
1. Installs `tailwindcss` and `@tailwindcss/vite`
2. Patches `vite.config.ts` — adds `import tailwindcss from '@tailwindcss/vite'` and wires `tailwindcss()` into the plugins array
3. Appends `@import "tailwindcss"` to `src/index.css`

### shadcn/ui
1. Installs `@types/node`
2. Patches `vite.config.ts` — adds `import path from 'path'` and a `resolve.alias` block mapping `@` → `./src`
3. Patches `tsconfig.app.json` and `tsconfig.json` — adds `baseUrl` and `paths` for the `@/*` alias
4. Runs `npx shadcn@latest init -d` (non-interactive defaults: New York style + CSS variables). If a non-default base color is selected in the wizard, appends `--base-color <theme>` (e.g. `slate`, `blue`, `rose`)

### ESLint + Prettier
1. Installs `eslint`, `prettier`, `eslint-config-prettier`, `@eslint/js`
2. Writes `eslint.config.js`, `.prettierrc`, `.prettierignore`
3. Adds `lint`, `lint:fix`, `format`, `format:check` scripts to `package.json`

### Biome
1. Installs `@biomejs/biome`
2. Runs `npx biome init`
3. Adds `lint`, `format`, `check` scripts

### Vitest
1. Installs `vitest` and `@vitest/ui`
2. Writes `vitest.config.ts`
3. Adds `test`, `test:ui`, `coverage` scripts

### Playwright
1. Installs `@playwright/test`
2. Writes `playwright.config.ts`
3. Runs `npx playwright install` (downloads browsers)
4. Adds `e2e`, `e2e:ui` scripts

### TanStack Query (React)
1. Installs `@tanstack/react-query`
2. Patches `src/main.tsx` — inserts `QueryClient` import and `<QueryClientProvider>` wrapper

### Pinia (Vue)
1. Installs `pinia`
2. Patches `src/main.ts` — inserts `createPinia()` call after `app` is created

### Angular Material
1. Runs `ng add @angular/material` (sets up theming, animations, gestures)

---

## Package Manager Support

All four package managers are supported. For new projects, pick at the wizard. For existing projects, the tool auto-detects from lockfiles:

| Lockfile | Package manager |
|----------|----------------|
| `pnpm-lock.yaml` | pnpm |
| `yarn.lock` | yarn |
| `bun.lockb` | bun |
| `package-lock.json` | npm |

Defaults to `npm` if no lockfile is found.

---

## Examples

### Scaffold a React SPA with shadcn/ui and TanStack Query

```
$ frontend-init

? Mode               New project
? Project name       my-dashboard
? Preset             Custom
? Package manager    pnpm
? Framework          react
? Variant            vite
? TypeScript         yes
? Linting            ESLint+Prettier
? UI library         shadcn/ui
? shadcn theme       zinc
? Testing            Vitest, Testing Library
? Tooling            TanStack Query, Zustand

┌─────────────────────────────────────────────┐
│  Review your setup                          │
│  Mode              new                      │
│  Package manager   pnpm                     │
│  Framework         react (vite)             │
│  TypeScript        true                     │
│  Linting           eslint-prettier          │
│  UI library        shadcn/ui                │
│  Testing           vitest, testing-library  │
│  Tooling           tanstack-query, zustand  │
└─────────────────────────────────────────────┘

Press Enter to execute...

✓  Scaffold project
✓  Install dependencies
✓  Install dev dependencies
✓  Configure ESLint + Prettier
⠹  Configure Tailwind CSS
   Patch files for shadcn/ui
   Update package.json scripts
   shadcn/ui: npx shadcn@latest init -d
   vitest: npx vitest init
```

### Use a preset (T3 Stack)

```
$ frontend-init

? Mode       New project
? Project name  my-t3-app
? Preset     T3 Stack

┌─────────────────────────────────────────────────────┐
│  Review your setup                                  │
│  Framework    react (nextjs)   UI     shadcn/ui     │
│  Linting      eslint-prettier  Test   vitest        │
│  Tooling      trpc, tanstack-query, zod             │
└─────────────────────────────────────────────────────┘
```

### Add tooling to an existing Vue project

```
$ cd my-existing-vue-app
$ frontend-init

? Mode       Configure existing directory
             (auto-detected: pnpm from pnpm-lock.yaml)

? Framework  vue
? Variant    vite
...
```

---

## Technical Notes

- **Angular CLI commands** (`ng new`, `ng add`) require `@angular/cli` to be installed globally: `npm install -g @angular/cli`
- **shadcn/ui** requires Tailwind CSS — the wizard automatically includes Tailwind when shadcn is selected
- **Existing project mode** skips the scaffold step, auto-detects the package manager from lockfiles, and uses the current directory as the project root
- **TypeScript** — all scaffold templates default to TypeScript variants when TypeScript is enabled

---

## Architecture

```
cmd/
  root.go          Cobra root command
  init.go          Launches wizard → executor pipeline
wizard/
  model.go         Root Bubbletea model + ExecuteModel
  flow.go          Dynamic step builder, framework-filtered options
  steps/
    step.go        Step interface
    select.go      Single-choice step (arrow keys)
    multiselect.go Checkbox multi-select step (space to toggle)
    input.go       Text input step (project name)
    confirm.go     Summary/review screen
    execute.go     Animated execution progress (spinner per task)
config/
  config.go        ProjectConfig struct
  presets.go       11 named preset definitions
executor/
  executor.go      Task list builder + orchestrator
  commands.go      Scaffold/install command builders
  configs.go       ToolSetup catalog + WriteConfigFiles
  patches.go       File patching (append, insert-after, replace)
  pkgjson.go       package.json script merger
  detect.go        Package manager auto-detection
```
