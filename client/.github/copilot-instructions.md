# Copilot Instructions for AI Coding Agents

## Project Overview
This is a React + TypeScript project bootstrapped with Vite. The codebase is organized for modularity and clarity, with a focus on component-driven development. The main entry point is `src/main.tsx`, and routing/pages are managed in `src/pages/`.

## Key Architecture & Patterns
- **Component Structure:**
  - UI components are grouped by feature in `src/components/` (e.g., `avatar`, `chat`, `emoji`, `inputchat`).
  - Pages are in `src/pages/` (e.g., `HomeChat.tsx`, `About.tsx`).
  - Assets (images, icons) are in `src/assets/`.
- **Styling:**
  - CSS files are colocated with components or in `src/` for global styles.
- **Utilities:**
  - Shared logic is in `src/utils/` (e.g., `date.ts`).
- **Type Safety:**
  - TypeScript is enforced throughout. Types are defined in `types.ts` files within component folders.

## Developer Workflows
- **Build:**
  - Use Vite for fast builds and HMR. Run `npm run dev` to start the development server.
- **Test:**
  - Jest is configured (see `jest.config.cjs`). Run `npm test` for unit tests. Coverage reports are generated in `coverage/`.
- **Linting:**
  - ESLint is set up with recommended and type-aware rules. See `eslint.config.js` for details. Use `npm run lint` to check code quality.
- **Type Checking:**
  - TypeScript configs are in `tsconfig.json`, `tsconfig.app.json`, and `tsconfig.node.json`.

## Conventions & Patterns
- **Component Exports:**
  - Each component folder has an `index.ts` for clean imports.
- **Naming:**
  - Use PascalCase for components, camelCase for functions/variables.
- **Testing:**
  - Place tests alongside utilities (e.g., `date.test.ts`).
- **Logging:**
  - Log files are stored in `logs/` for debugging API and UI flows.

## Integration Points
- **External Dependencies:**
  - React, Vite, TypeScript, Jest, ESLint, and plugins for React linting.
- **Cross-Component Communication:**
  - Props and context are used for data flow between components/pages.

## Examples
- To add a new UI feature, create a folder in `src/components/`, define your component and types, and export via `index.ts`.
- For a new page, add a `.tsx` file in `src/pages/` and update routing logic if present.

## References
- See `README.md` for setup and linting details.
- See `eslint.config.js` for custom lint rules.
- See `jest.config.cjs` for test configuration.

---

**If you are unsure about a pattern, check for similar usage in `src/components/` or `src/pages/`.**
