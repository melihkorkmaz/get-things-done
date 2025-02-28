# CLAUDE.md - GTD Application Guidelines

## Build & Run Commands
- Run application: `go run cmd/server/main.go`
- Install dependencies: `go mod tidy`
- Format code: `go fmt ./...`
- Lint code: `golint ./...` (install with: `go install golang.org/x/lint/golint@latest`)
- Test: `go test ./...`
- Test single package: `go test github.com/melihkorkmaz/gtd/internal/models`

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
- **Technology stack**: Go, Chi router, HTMX, Alpine.js, DaisyUI/Tailwind CSS