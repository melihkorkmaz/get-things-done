# CLAUDE.md - GTD Application Guidelines

## Build & Run Commands
- Run application: `go run cmd/server/main.go`
- Install dependencies: `go mod tidy`
- Format code: `go fmt ./...`
- Lint code: `golint ./...` (install with: `go install golang.org/x/lint/golint@latest`)
- Test: `go test ./...`
- Test single package: `go test github.com/melihkorkmaz/gtd/internal/models`
- Generate templ templates: `templ generate`

## Code Style Guidelines
- **Imports**: Group standard library, third-party packages, and internal packages
- **Error handling**: Use descriptive error messages, check errors immediately
- **Naming**: 
  - Variables/functions: camelCase
  - Exported functions/types: PascalCase
  - Packages: lowercase, single word
- **Comments**: Document all exported functions, types, and methods
- **File structure**: One package per directory, main.go for executables
- **HTML/CSS**: Follow BEM naming convention for CSS classes
- **Technology stack**: Go, Chi router, Templ, HTMX, Alpine.js, DaisyUI/Tailwind CSS, PostgreSQL

## Templ Guidelines
- Use `templ.Script` for JavaScript function calls in components
- Remember to make helper functions exportable (PascalCase) if used across components
- For HTMX integration, use attributes like `hx-get`, `hx-post`, etc.
- Render components with component.Render() in HTTP handlers
- Generate code with `templ generate` after making changes to .templ files

## Project Structure
- **/cmd/server**: Main application entry point
- **/internal/config**: Configuration handling (database, environment)
- **/internal/handlers**: HTTP handlers organized by feature
- **/internal/models**: Domain models and data store interfaces
- **/internal/views**: Templ template components
  - **/layouts**: Base page layouts (base.templ)
  - **/pages**: Full page templates (projects.templ, tasks_list.templ)
  - **/partials**: Reusable component templates (task_row.templ, project_card.templ)
- **/static**: Static assets (CSS, JavaScript)