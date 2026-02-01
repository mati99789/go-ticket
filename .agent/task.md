# Mentoring Plan: Project "GoTicket"

## Globalne Cele

- [x] Opanowanie zaawansowanego Go (Concurrency, Interfaces, Reflect).
- [ ] Zrozumienie architektury mikroserwisowej i komunikacji gRPC.
- [x] Praktyczne zastosowanie wzorców projektowych i DDD/Hexagonal.
- [ ] DevOps: Docker, CI/CD, AWS, Monitoring.

## Phase 0: Setup & Foundation

- [x] Wybór i zatwierdzenie tematu projektu. <!-- id: 0 -->
- [x] Konfiguracja Repozytorium (Monorepo). <!-- id: 1 -->
- [x] Setup środowiska lokalnego (Go, Docker, Makefiles). <!-- id: 2 -->
- [x] Projekt architektury wysokopoziomowej (C4 Model). <!-- id: 3 -->

## Phase 1: The Domain & Core

- [x] Struktura projektu (Go Standard Layout). <!-- id: 4 -->
- [x] Implementacja domeny `Events` (Structs, Entities). <!-- id: 5 -->
- [x] Implementacja domeny `Booking` (Statuses, Validation). <!-- id: 18 -->
- [x] Implementacja warstwy dostępu do danych (Postgres + pgx/sqlc). <!-- id: 6 -->

## Phase 2: API, Middleware & Transactions

- [x] Wystawienie REST API (net/http). <!-- id: 7 -->
- [x] Middleware (Logging, Recovery). <!-- id: 8 -->
- [x] Implementacja logicznej transakcji (Service Layer). <!-- id: 19 -->
- [ ] Wprowadzenie gRPC (Proto definitions). <!-- id: 9 -->

## Phase 3: The Hard Parts (Concurrency & QA)

- [x] System rezerwacji (Booking Logic + Atomic Reservation). <!-- id: 10 -->
- [ ] Testy obciążeniowe (k6) - sprawdzenie Race Conditions. <!-- id: 12 -->
- [ ] Implementacja blokowania zasobów (Redis Distributed Locks - EARN THE OVERKILL). <!-- id: 11 -->

## Phase 4: Polish & Production Standards

- [x] Refactor: DTO Layer (czyste odpowiedzi JSON). <!-- id: 20 -->
- [x] Refactor: Advanced Error Mapping (Domain Errors). <!-- id: 21 -->
- [ ] **Pagination & Query Parameters**: <!-- id: 37 -->
  - [ ] Implement pagination (limit, offset) for ListEvents. <!-- id: 38 -->
  - [ ] Add filtering (by date range, price range). <!-- id: 39 -->
  - [ ] Add sorting (by name, price, date). <!-- id: 40 -->
  - [ ] Cursor-based pagination (for high performance). <!-- id: 41 -->
- [ ] **Authentication & Authorization**: <!-- id: 30 -->
  - [ ] User Registration & Login (JWT). <!-- id: 31 -->
  - [ ] Password hashing (bcrypt). <!-- id: 32 -->
  - [ ] JWT Middleware (protected routes). <!-- id: 33 -->
  - [ ] RBAC implementation (admin/organizer/user roles). <!-- id: 34 -->
  - [ ] OAuth2 integration (Google/GitHub). <!-- id: 35 -->
  - [ ] Refresh Token mechanism (Redis). <!-- id: 36 -->
- [x] Automatyzacja migracji w kodzie (Embedded migrations). <!-- id: 17 -->
- [/] Testy integracyjne (Database tests). <!-- id: 22 -->
  - [x] Test Infrastructure (Testcontainers, fixtures)
  - [x] EventRepository tests (CRUD, race conditions)
  - [ ] BookingService tests (transactions)
  - [ ] Handler tests (end-to-end API)

## Phase 5: DevOps & Cloud

- [ ] Konteneryzacja (Dockerfile, Docker Compose dla całości). <!-- id: 13 -->
- [ ] **Observability & Monitoring**: <!-- id: 42 -->
  - [ ] Prometheus metrics (HTTP, business metrics). <!-- id: 43 -->
  - [ ] Grafana dashboards (API performance, business KPIs). <!-- id: 44 -->
  - [ ] Structured logging with slog (JSON output). <!-- id: 45 -->
  - [ ] Distributed tracing (OpenTelemetry). <!-- id: 46 -->
  - [ ] Health check endpoints (/health, /ready). <!-- id: 47 -->
- [ ] GitHub Actions (CI pipelines). <!-- id: 14 -->
- [ ] Provisioning AWS (Terraform/OpenTofu). <!-- id: 15 -->
- [ ] Deploy na AWS. <!-- id: 16 -->

## Phase 6: Senior Upgrades (The "Antigravity" Audit)

- [ ] **Tests**: Integration Tests (Testcontainers) for EventRepository. <!-- id: 23 -->
- [ ] **Docs**: Swagger/OpenAPI documentation. <!-- id: 24 -->
- [ ] **API**: GraphQL Layer (gqlgen). <!-- id: 25 -->
- [ ] **Cache**: Redis implementation. <!-- id: 26 -->
- [ ] **Async**: RabbitMQ/Kafka for Domain Events. <!-- id: 27 -->
- [ ] **Prod**: Multi-stage Dockerfile (Distroless). <!-- id: 28 -->
- [ ] **Automation**: Makefile. <!-- id: 29 -->
