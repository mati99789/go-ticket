# Migration Summary: Antigravity IDE → Claude Code

## What Was Done

Migrated all project configuration and mentorship files from Antigravity IDE format to Claude Code format. Converted all Polish-language content to English. Restructured the entry point to match Claude Code's architecture.

---

## Files Changed

### Created

| File | Why |
|------|-----|
| `CLAUDE.md` | Claude Code's primary auto-loaded instruction file. Contains the full mentorship philosophy (Iron Law, Socratic method, red flags, DoD summary, resource references). This replaces `.agent/rules.md` as the entry point Claude Code reads automatically. |
| `.claude/commands/restore-context.md` | Proper Claude Code slash command (`/restore-context`). Replaces `.agent/workflows/restore_context.md` which used Antigravity-specific `// turbo` syntax that Claude Code does not understand. |

---

### Modified

| File | What Changed | Why |
|------|--------------|-----|
| `.agent/rules.md` | Full rewrite: Polish → English. Removed Go code blocks (violated the "no ready code" rule). Kept all mentor rules, Red Flags section converted from code examples to descriptive questions. | 37 Polish character occurrences. Code blocks in a "no code" mentor document was a contradiction. |
| `.agent/definition_of_done.md` | Full rewrite: Polish → English. Added Race Clean check (`go test -race`), ADR requirement, and Roadmap update requirement. | 10 Polish character occurrences. Missing critical checks. |
| `.agent/git_conventions.md` | Full rewrite: Polish → English. Added branching strategy section and "Why does this matter?" explanation. | 4 Polish character occurrences. Missing branch naming conventions. |
| `.agent/tech_stack.md` | Full rewrite: Polish → English. Restructured into tables for readability. Added Kafka + RabbitMQ (both were missing from original). Added "Why these choices?" rationale. | 6 Polish character occurrences. Kafka/RabbitMQ already implemented but not in the approved stack doc. |
| `.agent/c4_architecture.md` | Converted Polish comments to English. Updated Kafka reference (was RabbitMQ in the microservices diagram, now correctly reflects the actual implementation). Added Level 3 component diagram for Booking module. | 3 Polish character occurrences. Mermaid diagram was inaccurate (showed RabbitMQ where Kafka is used). |
| `.agent/project_roadmap.md` | Converted Polish text to English. Updated Phase 3 to reflect actual completion status of Kafka SyncProducer. Added Phase 4 (Notification Service) with concrete sub-tasks. | 4 Polish character occurrences. Phase completion statuses were out of date. |
| `.agent/task.md` | Converted all Polish inline text to English. Fixed "Security Design Issue" section to describe the fix clearly. | 8 Polish character occurrences. Mixed language made the document hard to read. |
| `.agent/workflows/restore_context.md` | Removed Antigravity `// turbo` syntax. Converted to standard markdown with a pointer to the new `/restore-context` slash command. | Antigravity-specific syntax breaks in Claude Code. |
| `tests/load/README.md` | Full rewrite: Polish → English. Added "Proven Result" section documenting the verified 10000 = bookings + remaining spots result. Added scenario structure documentation. | 13 Polish character occurrences. Missing proof of test results. |

---

### Not Changed (Already Good)

| File | Reason |
|------|--------|
| `README.md` | Already in English, comprehensive, accurate |
| `PROJECT_ROADMAP.md` | Already in English, up to date |
| `TECH_STACK.md` | Already in English (0 Polish characters found) |
| `make-tutorial-command.md` | Already in English, accurate |
| `.agent/adr_template.md` | Already in English, correct format |
| `.agent/project_plan.md` | Already in English, accurate architecture description |

---

## Why These Changes Were Necessary

1. **Claude Code reads `CLAUDE.md` automatically** — without it, the mentorship rules were invisible to Claude Code at session start.

2. **`.claude/commands/` is how Claude Code handles slash commands** — the Antigravity `// turbo` directive has no meaning in Claude Code and the workflow would never execute.

3. **Polish language in mentor files** — Claude Code works in English by default. Mixed-language files produced inconsistent mentor behavior and were inaccessible to non-Polish speakers reviewing the project.

4. **Code blocks in `rules.md`** — the Iron Law says "zero ready code". Having Go code blocks in the mentor rules file was a direct contradiction of the mentorship philosophy, even if they were showing anti-patterns.

5. **Stale status in roadmap files** — the project_roadmap showed RabbitMQ in the microservices diagram, but the actual implementation uses Kafka. Documentation that contradicts the code is worse than no documentation.
