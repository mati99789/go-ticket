# Git & Commit Conventions

A professional project requires a clean, readable change history. We follow **Conventional Commits**.

---

## Format

```
<type>(<scope>): <subject>
```

---

## Types

| Type | When to use |
|------|-------------|
| `feat` | New functionality (new endpoint, new struct, new feature) |
| `fix` | Bug fix |
| `docs` | Documentation changes only |
| `style` | Formatting changes, no logic changes (gofmt, lint) |
| `refactor` | Code change without behavior change (optimization, cleanup) |
| `test` | Adding or improving tests |
| `chore` | Technical tasks (library updates, build config, CI/CD) |

---

## Examples

```
feat(domain): add Event struct with capacity validation
fix(api): handle timeout error in booking endpoint
docs(agent): update mentorship rules to English
test(event): add table driven tests for price calculation
refactor(main): extract setupRoutes and gracefulShutdown helpers
chore(ci): upgrade golangci-lint to v2.10.1
```

---

## Rules

1. **Imperative Mood**: Write "add", not "added" or "adds"
2. **English Only**: All commit messages must be in English — no exceptions
3. **Atomic Commits**: One commit = one logical change. Do not bundle unrelated changes
4. **No WIP commits**: Never commit with message "wip", "temp", "fix typo". Squash before merge
5. **Scope is optional but recommended**: Use the affected package/module as scope

---

## Branching Strategy

```
main          — production-ready, protected
feature/*     — new features (e.g., feature/kafka-consumer)
fix/*         — bug fixes (e.g., fix/booking-race-condition)
chore/*       — maintenance (e.g., chore/update-dependencies)
```

---

> **Why does this matter?** A clean git history is a communication tool. In 6 months, when a bug appears in production, a well-written commit history tells the story of what changed, when, and why — without needing to read every line of code.
