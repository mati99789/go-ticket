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
  - [ ] Handler tests (end-to-end API Mocking)
- [ ] Load Testing (k6) - Verify Race Conditions under load. <!-- id: 12 -->

## Phase 4: Security & Advanced Logic

- [ ] **Authentication & Authorization**: <!-- id: 30 -->
  - [ ] User Registration & Login (JWT). <!-- id: 31 -->
  - [ ] Password hashing (bcrypt). <!-- id: 32 -->
  - [ ] JWT Middleware (protected routes). <!-- id: 33 -->
  - [ ] RBAC implementation (admin/organizer/user roles). <!-- id: 34 -->
- [ ] **Pagination & Query Parameters**: <!-- id: 37 -->
  - [ ] Implement pagination (limit, offset) for ListEvents. <!-- id: 38 -->
  - [ ] Add filtering (by date range, price range). <!-- id: 39 -->

## Phase 5: DevOps & Cloud

- [ ] Containerization (Dockerfile, Docker Compose). <!-- id: 13 -->
- [ ] GitHub Actions (CI pipelines). <!-- id: 14 -->
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
