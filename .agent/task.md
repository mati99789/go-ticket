# Mentoring Plan: Project "GoTicket"

## Global Goals

- [x] Master Advanced Go (Concurrency, Interfaces, Reflect).
- [ ] Understand Microservices Architecture and gRPC communication.
- [x] Practical application of Design Patterns & DDD/Hexagonal.
- [ ] DevOps: Docker, CI/CD, AWS, Monitoring, Kubernetes.

---

## Phase 0: Setup & Foundation

- [x] Project Topic Selection & Approval. <!-- id: 0 -->
- [x] Repository Configuration (Monorepo). <!-- id: 1 -->
- [x] Local Environment Setup (Go, Docker, Makefiles). <!-- id: 2 -->
- [x] High-Level Architecture Design (C4 Model). <!-- id: 3 -->

---

## Phase 1: The Domain & Core

- [x] Project Structure (Go Standard Layout). <!-- id: 4 -->
- [x] Domain Implementation `Events` (Structs, Entities). <!-- id: 5 -->
- [x] Domain Implementation `Booking` (Statuses, Validation). <!-- id: 18 -->
- [x] Data Access Layer Implementation (Postgres + pgx/sqlc). <!-- id: 6 -->

---

## Phase 2: API, Middleware & Transactions

- [x] REST API Implementation (net/http). <!-- id: 7 -->
- [x] Middleware (Logging, Recovery). <!-- id: 8 -->
- [x] Service Layer Transaction Implementation. <!-- id: 19 -->
- [x] DTO Layer Refactor (Clean JSON responses). <!-- id: 20 -->
- [x] Advanced Error Mapping Refactor (Domain Errors). <!-- id: 21 -->
- [x] Embedded Migrations Automation. <!-- id: 17 -->

---

## Phase 3: The Hard Parts (Concurrency & QA)

- [x] Reservation System (Booking Logic + Atomic Reservation). <!-- id: 10 -->
- [x] Integration Tests (Database tests). <!-- id: 22 -->
  - [x] Test Infrastructure (Testcontainers, fixtures)
  - [x] EventRepository tests (CRUD, race conditions)
  - [x] BookingService tests (transactions, rollback)
  - [x] Handler tests (end-to-end API mocking)
- [x] Load Testing (k6) — Verify Race Conditions under load. <!-- id: 12 -->
  - [x] `tests/load/seed.sql` — organizer user + event (capacity=10000), ON CONFLICT DO NOTHING
  - [x] `tests/load/booking_scenario.js` — setup() register+login, default(data) with token and eventId
  - [x] `tests/load/README.md` — full documentation: how to run, seed, verify after test, CI/CD plan
  - [x] Run seed: `docker exec -i go_ticket_db psql -U postgres -d go_ticket < tests/load/seed.sql`
  - [x] Run test: `k6 run tests/load/booking_scenario.js`
  - [x] Verify result: 1031 + 8969 = 10000 — zero double-bookings confirmed

---

## Phase 4: Security & Advanced Logic

- [x] **Authentication & Authorization**: <!-- id: 31 -->
  - [x] User Registration & Login (JWT). <!-- id: 30 -->
    - [x] User Domain Entity (`internal/domain/user.go`) — validation, roles (user/admin/organizer)
    - [x] Database Migration (`000003_add_user_table.up.sql`) — users table with ENUM roles
    - [x] Password Hashing (`internal/auth/password.go`) — bcrypt implementation
    - [x] JWT Token Generation (`internal/auth/jwt.go`) — GenerateToken, VerifyToken
    - [x] SQLC Queries (`internal/postgres/queries/users.sql`) — CreateUser, GetUserByEmail, etc.
    - [x] User Repository (`internal/postgres/user_repository.go`)
    - [x] Auth Handler (`internal/api/auth_handler.go`) — /register, /login endpoints
    - [x] Wire Auth endpoints in main.go
  - [x] JWT Middleware (`internal/api/middleware/auth.go`) — protected routes. <!-- id: 33 -->
  - [x] RBAC implementation (`RequireRole()` middleware + `requireOrganizer`/`requireAdmin`/`requireAll` wrappers in main.go). <!-- id: 34 -->
- [x] **Security Middleware**: <!-- id: 35 -->
  - [x] Rate Limiting Middleware — Redis INCR+EXPIRE, IPKey (auth) + UserKey (API), miniredis tests. <!-- id: 36 -->
  - [ ] Audit Logging Middleware (compliance). <!-- id: 37 -->
- [x] **Security Fix — Booking userEmail**: <!-- id: 68 -->
  - [x] `CreateBooking` now extracts userEmail from JWT claims, not request body
  - [x] Prevents authenticated users from creating bookings on behalf of other emails
- [ ] **Pagination & Query Parameters**: <!-- id: 38 -->
  - [ ] Implement pagination (limit, offset) for ListEvents. <!-- id: 39 -->
  - [ ] Add filtering (by date range, price range). <!-- id: 40 -->

---

## Phase 5: DevOps & Cloud

- [x] Containerization (Dockerfile multi-stage distroless, Docker Compose app+db). <!-- id: 13 -->
- [x] GitHub Actions (CI + CD pipelines, golangci-lint v2, docker push to GHCR). <!-- id: 14 -->
- [ ] Provisioning AWS (Terraform/OpenTofu). <!-- id: 15 -->
- [ ] **Kubernetes (K8s)**: <!-- id: 48 -->
  - [ ] Deploy local cluster (Kind/Minikube). <!-- id: 49 -->
  - [ ] Helm Charts for GoTicket. <!-- id: 50 -->

---

## Phase 6: Senior Upgrades — Event Streaming

- [x] Swagger/OpenAPI Documentation <!-- id: 70 -->
  - [x] Install and configure `swag` and `http-swagger`
  - [x] Expose Swagger UI at `/swagger/`
  - [x] Fix annotations to use DTO structs instead of loose body parameters
- [ ] Observability (Prometheus, Grafana, Distributed Tracing — Jaeger/OpenTelemetry) <!-- id: 50 -->
- [ ] **gRPC**: Proto definitions and service implementation. <!-- id: 9 -->
- [ ] **Cache**: Redis caching layer (beyond rate limiting). <!-- id: 11 -->
- [/] **Async**: Kafka / RabbitMQ for Domain Events. <!-- id: 27 -->
  - [x] Implement Transactional Outbox Pattern in PostgreSQL (prevents dual-write problem)
  - [x] Implement KafkaBroker SyncProducer (`internal/kafka/kafka_broker.go` via IBM/sarama)
  - [x] Wire Kafka Broker into `cmd/app/main.go` via `setupKafkaRelay`
  - [x] Implement Kafka Consumer for `booking_events_topic`
    - [x] `domain.MessageConsumer` interface (`internal/domain/message_consumer.go`)
    - [x] `domain.EventHandler` interface (`internal/domain/eventhandle.go`)
    - [x] `KafkaConsumer` implementation (`internal/kafka/kafka_consumer.go`)
    - [x] `KafkaConsumerWorker` with exponential backoff + interruptible sleep (`internal/workers/kafka_consumer_worker.go`)
    - [x] `BookingEventHandler` (`internal/event_handler/booking_event_handler.go`)
    - [x] Wire consumer in `main.go` via `setupKafkaConsumer` with graceful shutdown
    - [x] Error handling via `erChan` (buffer 3) for relay, consumer, server goroutines
    - [x] Fixed shadow variable `ctx` → `initialCtx` / `workerCtx`
  - [ ] Implement RabbitMQ Publisher for email task queue
  - [ ] Implement Email Worker goroutine
- [ ] **Microservices**: Extract Notification/Payment Service. <!-- id: 28 -->

---

## Phase 7: Production Polish

> **Goal**: Add critical production features with high impact and low effort.
> **Estimated Time**: 2–3 weeks (part-time)

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
