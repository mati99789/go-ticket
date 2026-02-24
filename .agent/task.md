# Mentoring Plan: Project "GoTicket"

## Global Goals

- [x] Master Advanced Go (Concurrency, Interfaces, Reflect).
- [ ] Understand Microservices Architecture and gRPC communication.
- [x] Practical application of Design Patterns & DDD/Hexagonal.
- [ ] DevOps: Docker, CI/CD, AWS, Monitoring, Kubernetes.

## Phase 0: Setup & Foundation

- [x] Project Topic Selection & Approval. <!-- id: 0 -->
- [x] Repository Configuration (Monorepo). <!-- id: 1 -->
- [x] Local Environment Setup (Go, Docker, Makefiles). <!-- id: 2 -->
- [x] High-Level Architecture Design (C4 Model). <!-- id: 3 -->

## Phase 1: The Domain & Core

- [x] Project Structure (Go Standard Layout). <!-- id: 4 -->
- [x] Domain Implementation `Events` (Structs, Entities). <!-- id: 5 -->
- [x] Domain Implementation `Booking` (Statuses, Validation). <!-- id: 18 -->
- [x] Data Access Layer Implementation (Postgres + pgx/sqlc). <!-- id: 6 -->

## Phase 2: API, Middleware & Transactions

- [x] REST API Implementation (net/http). <!-- id: 7 -->
- [x] Middleware (Logging, Recovery). <!-- id: 8 -->
- [x] Service Layer Transaction Implementation. <!-- id: 19 -->
- [x] DTO Layer Refactor (Clean JSON responses). <!-- id: 20 -->
- [x] Advanced Error Mapping Refactor (Domain Errors). <!-- id: 21 -->
- [x] Embedded Migrations Automation. <!-- id: 17 -->

## Phase 3: The Hard Parts (Concurrency & QA)

- [x] Reservation System (Booking Logic + Atomic Reservation). <!-- id: 10 -->
- [x] Integration Tests (Database tests). <!-- id: 22 -->
  - [x] Test Infrastructure (Testcontainers, fixtures)
  - [x] EventRepository tests (CRUD, race conditions)
  - [x] BookingService tests (transactions, rollback)
  - [x] Handler tests (end-to-end API Mocking)
- [ ] Load Testing (k6) - Verify Race Conditions under load. <!-- id: 12 -->

## Phase 4: Security & Advanced Logic

- [x] **Authentication & Authorization**: <!-- id: 31 -->
  - [x] User Registration & Login (JWT). <!-- id: 30 -->
    - [x] User Domain Entity (`internal/domain/user.go`) - validation, roles (user/admin/organizer)
    - [x] Database Migration (`000003_add_user_table.up.sql`) - users table with ENUM roles
    - [x] Password Hashing (`internal/auth/password.go`) - bcrypt implementation
    - [x] JWT Token Generation (`internal/auth/jwt.go`) - GenerateToken, VerifyToken
    - [x] SQLC Queries (`internal/postgres/queries/users.sql`) - CreateUser, GetUserByEmail, etc.
    - [x] User Repository (`internal/postgres/user_repository.go`)
    - [x] Auth Handler (`internal/api/auth_handler.go`) - /register, /login endpoints
    - [x] Wire Auth endpoints in main.go
  - [x] JWT Middleware (`internal/api/middleware/auth.go`) - protected routes. <!-- id: 33 -->
  - [x] RBAC implementation (`RequireRole()` middleware + `requireOrganizer`/`requireAdmin`/`requireAll` wrappers in main.go). <!-- id: 34 -->
- [ ] **Security Middleware**: <!-- id: 35 -->
  - [ ] Rate Limiting Middleware (anti-scanning protection). <!-- id: 36 -->
  - [ ] Audit Logging Middleware (compliance). <!-- id: 37 -->
- [ ] **Pagination & Query Parameters**: <!-- id: 38 -->
  - [ ] Implement pagination (limit, offset) for ListEvents. <!-- id: 39 -->
  - [ ] Add filtering (by date range, price range). <!-- id: 40 -->

## Phase 5: DevOps & Cloud

- [x] Containerization (Dockerfile multi-stage distroless, Docker Compose app+db). <!-- id: 13 -->
- [x] GitHub Actions (CI + CD pipelines, golangci-lint v2, docker push to GHCR). <!-- id: 14 -->
- [ ] Provisioning AWS (Terraform/OpenTofu). <!-- id: 15 -->
- [ ] **Kubernetes (K8s)**: <!-- id: 48 -->
  - [ ] Deploy local cluster (Kind/Minikube). <!-- id: 49 -->
  - [ ] Helm Charts for GoTicket. <!-- id: 50 -->
- [ ] Observability & Monitoring (Prometheus, Grafana, OpenTelemetry). <!-- id: 42 -->

## Phase 6: Senior Upgrades (The "Antigravity" Audit)

- [ ] **RPC**: gRPC Implementation (Proto definitions). <!-- id: 9 -->
- [ ] **Docs**: Swagger/OpenAPI documentation. <!-- id: 24 -->
- [ ] **Cache**: Redis implementation (Caching & Distributed Locks). <!-- id: 11 -->
- [ ] **Async**: Kafka/RabbitMQ for Domain Events. <!-- id: 27 -->
- [ ] **Microservices**: Extracting Notification/Payment Service. <!-- id: 28 -->

## Phase 7: Production Polish (High-Impact Additions)

> **Goal:** Add critical production features with high impact and low effort.
> **Estimated Time:** 2-3 weeks (part-time)

- [ ] **Health Checks**: <!-- id: 51 -->
  - [ ] `/health` endpoint (basic health check). <!-- id: 52 -->
  - [ ] `/readiness` endpoint (K8s readiness probe). <!-- id: 53 -->
  - [ ] `/liveness` endpoint (K8s liveness probe). <!-- id: 54 -->
- [ ] **HTTPS/TLS**: <!-- id: 55 -->
  - [ ] cert-manager setup in K8s. <!-- id: 56 -->
  - [ ] Let's Encrypt integration. <!-- id: 57 -->
- [ ] **Secrets Management**: <!-- id: 58 -->
  - [ ] K8s Secrets configuration. <!-- id: 59 -->
  - [ ] Sealed Secrets (GitOps-friendly). <!-- id: 60 -->
- [ ] **Error Tracking**: <!-- id: 61 -->
  - [ ] Sentry integration (free tier). <!-- id: 62 -->
  - [ ] Error reporting middleware. <!-- id: 63 -->
- [ ] **Resilience Patterns**: <!-- id: 64 -->
  - [ ] Circuit Breaker implementation (go-resilience). <!-- id: 65 -->
  - [ ] Retry Logic with exponential backoff. <!-- id: 66 -->
  - [ ] Timeout configuration for external calls. <!-- id: 67 -->
