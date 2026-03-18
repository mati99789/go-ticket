# Definition of Done (DoD)

For a task (Ticket/Task) to be considered complete, it must meet **all** of the following criteria. Missing any single point = Rejected in Code Review.

---

## 1. Code Quality

- [ ] **Linter Clean**: `golangci-lint run` reports zero errors
- [ ] **Zero Magic Numbers/Strings**: All constants are extracted to `const` with meaningful names
- [ ] **Meaningful Names**: Variables `x`, `data`, `temp` are forbidden — names must reflect business intent
- [ ] **Error Handling**: Every error is handled or returned (wrapped with context). Ignoring errors (`_`) is forbidden

---

## 2. Tests

- [ ] **Unit Tests**: Business logic has test coverage using Table-Driven Tests
- [ ] **Green Build**: All tests (`go test ./...`) pass with zero failures
- [ ] **Race Clean**: `go test -race ./...` passes with zero race conditions detected

---

## 3. Architecture

- [ ] **Dependency Rule**: The domain layer does not import external libraries (Database, HTTP, etc.)
- [ ] **Single Responsibility**: Every function/method does exactly one thing — if you need "and" to describe it, split it

---

## 4. Git & History

- [ ] **Commit Message**: Follows Conventional Commits format (e.g., `feat: add event structure`, `fix: calculation error`)
- [ ] **Clean History**: No "wip", "fix typo", or "temp" commits. Squash before merge
- [ ] **English Only**: All commit messages and code comments are in English

---

## 5. Documentation (for significant changes)

- [ ] **ADR**: If an architectural decision was made, an ADR is created using `.agent/adr_template.md`
- [ ] **Roadmap updated**: `.agent/task.md` reflects the completed task

---

> **Mentorship Note**: The DoD is not a bureaucratic checklist — it is the standard that separates "it works on my machine" from "it works in production."
