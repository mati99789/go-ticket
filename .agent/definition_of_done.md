# Definition of Done (DoD)

Aby zadanie (Ticket/Task) zostało uznane za zakończone, musi spełniać poniższe kryteria. Brak któregokolwiek punktu = Odrzucenie w Code Review.

## 1. Jakość Kodu
*   [ ] **Linter Clean**: `golangci-lint run` nie zgłasza żadnych błędów.
*   [ ] **Zero Magic Numbers/Strings**: Wszystkie stałe są wyniesione do `const`.
*   [ ] **Meaningful Names**: Zmienne `x`, `data`, `temp` są zakazane. Nazwa musi oddawać intencję biznesową.
*   [ ] **Error Handling**: Każdy błąd jest obsłużony (handle) lub zwrócony (wrap with context). Ignorowanie błędów (`_`) jest zabronione.

## 2. Testy
*   [ ] **Unit Tests**: Logika biznesowa ma pokrycie testami (Table Driven Tests).
*   [ ] **Green Build**: Wszystkie testy (`go test ./...`) przechodzą.

## 3. Architektura
*   [ ] **Dependency Rule**: Domena nie zależy od zewnętrznych bibliotek (Database, HTTP).
*   [ ] **Responsibility**: Funkcja/Metoda robi tylko jedną rzecz (Single Responsibility Principle).

## 4. Git & Historia
*   [ ] **Commit Message**: Zgodny z Conventional Commits (np. `feat: add event structure`, `fix: calculation error`).
*   [ ] **Clean History**: Brak commitów typu "wip", "fix tyop". Squash przed merge'em.
