# Git & Commit Conventions

Profesjonalny projekt wymaga czystej historii zmian. Stosujemy **Conventional Commits**.

## Format
`<type>(<scope>): <subject>`

### Types
*   `feat`: Nowa funkcjonalność (nowy endpoint, nowa struktura).
*   `fix`: Naprawa błędu.
*   `docs`: Zmiany w dokumentacji.
*   `style`: Formatowanie, brak zmian w logice (gofmt, lint).
*   `refactor`: Zmiana kodu bez zmiany zachowania (optymalizacja, czyszczenie).
*   `test`: Dodanie lub poprawa testów.
*   `chore`: Zadania techniczne (update bibliotek, konfig builda).

### Przykłady
*   `feat(domain): add Event struct definition`
*   `fix(api): handle timeout in booking endpoint`
*   `docs(agent): update mentorship rules`
*   `test(event): add table driven tests for price calculation`

## Zasady
1.  **Imperative Mood**: Pisz "add", a nie "added" czy "adds".
2.  **English Only**: Komunikaty piszemy po angielsku.
3.  **Atomic Commits**: Jeden commit = jedna logiczna zmiana.
