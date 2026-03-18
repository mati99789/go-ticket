# Mentorship Roadmap: go-ticket Project

This document tracks the current state of the `go-ticket` project and the roadmap to transform it into a Senior-Level Backend/DevOps portfolio project.

---

## 1. Project Audit: Current State vs. Target

| Category | Component | Current State | Target "Senior" State | Status |
| :--- | :--- | :--- | :--- | :--- |
| **Architecture** | Pattern | Modular Monolith (Clean Architecture) | DDD with clear bounded context boundaries | 🟢 Solid |
| | Separation | Domain / Service / Repo / API | Hexagonal / Clean Architecture (Strict) | 🟢 Done |
| **Backend (Go)** | API | REST (Standard Lib) | REST + GraphQL + gRPC | 🟡 REST done |
| | Concurrency | Atomic booking, Testcontainers | Advanced Patterns (Workers, Pipelines) | 🟢 Good |
| | Persistence | Postgres (pgx + sqlc) | Postgres + Redis (Caching + Rate Limiting) | 🟢 Done |
| **Security** | Authentication | JWT (access token) | JWT + OAuth2 (Google/GitHub) | 🟢 JWT done |
| | Authorization | RBAC (user/organizer/admin) | Full RBAC + fine-grained permissions | 🟢 Done |
| **DevOps** | Containerization | Multi-stage Dockerfile (distroless) | Optimized Multi-stage Dockerfiles | 🟢 Done |
| | Orchestration | None | Kubernetes (Helm/Kustomize) | 🔴 Planned |
| | IaC | None | Terraform / OpenTofu | 🔴 Planned |
| | CI/CD | GitHub Actions (CI + CD) | GitHub Actions (Lint, Test, Build, Push) | 🟢 Done |
| **Messaging** | Async | Outbox + Kafka Producer + Consumer | RabbitMQ (Task Queue) + Email Worker | 🟡 Kafka done |
| **Quality** | Testing | Integration Tests (Testcontainers) | E2E, Load (k6), Property-based Tests | 🟢 Int. tests done |
| | Observability | Structured logging (slog/JSON) | Distributed Tracing (OTEL), Metrics (Prometheus) | 🟡 Logging only |

---

## 2. Evolutionary Architecture Roadmap

We do not build everything at once. We follow an iterative "Evolutionary Architecture" approach.

### Phase 1: Foundation Hardening ✅ COMPLETE

- [x] Testing Strategy: Integration Tests (Testcontainers) for Repositories + API Handlers
- [x] Structured Logging: JSON logging via `slog`
- [x] Security & Auth: JWT, RBAC (user/organizer/admin), bcrypt, user enumeration protection
- [x] Error Handling: Domain errors → HTTP status codes mapping
- [x] Embedded Migrations: Auto-run on startup

### Phase 2: Security & DevOps Foundations ✅ COMPLETE

- [x] Dockerfile: Multi-stage build (builder → distroless/static-debian12), 31MB image
- [x] docker-compose: app + postgres services, healthcheck, env_file
- [x] GitHub Actions CI: lint (golangci-lint v2) + test (-race, coverage) + build on every push
- [x] GitHub Actions CD: docker build + push to GHCR on main (workflow_run)
- [x] Code Quality: Fixed 36 golangci-lint issues (errcheck, gosec, lll, gocyclo, funlen)
- [x] Rate Limiting: Redis INCR+EXPIRE, two limiters (auth: 5/15min IP-based, API: 100/1min user-based), miniredis unit tests
- [x] main.go refactor: extracted setupRoutes, setupServer, gracefulShutdown, setupRepositories
- [x] MapDomainError refactor: table-driven lookup with errorMapping struct
- [x] Load Testing (k6): verified race-condition safety under load
  - [x] `tests/load/seed.sql` — organizer user + event with capacity=10000
  - [x] `tests/load/booking_scenario.js` — setup() login→token + default(data) with auth header
  - [x] `tests/load/README.md` — full documentation: how to run, seed, verify, CI/CD plan
  - [x] Verified: 1031 booked + 8969 remaining = 10000 — zero double-bookings confirmed

### Phase 3: Event Streaming (Microservices Prep) 🔄 IN PROGRESS

- [x] Transactional Outbox Pattern in PostgreSQL (prevents dual-write problem)
- [x] OutboxRelay background worker (Goroutine loop with graceful shutdown)
- [x] Kafka KRaft Docker infrastructure (INTERNAL:9092 + EXTERNAL:9094 listeners)
- [x] `domain.MessageBroker` interface + Kafka `SyncProducer` implementation (IBM/sarama)
- [x] Wire Kafka Broker into `cmd/app/main.go` via `setupKafkaRelay`
- [x] **Swagger/OpenAPI**: Auto-generated docs (`swaggo`) — Swagger UI at `/swagger/`
- [ ] **Correlation IDs**: Request tracing through middleware — essential for distributed debugging
- [ ] **Pagination & Filtering**: `limit/offset` for ListEvents + date/price range filters

### Phase 4: Notification Service (RabbitMQ + Email) 🔄 IN PROGRESS

> Prerequisites: Complete Phase 3

- [x] Kafka Consumer for `booking_events_topic`
- [x] Wire consumer goroutine in `main.go` with graceful shutdown + `erChan` error propagation
- [ ] RabbitMQ added to `docker-compose.yml`
- [ ] RabbitMQ Publisher (`internal/rabbitmq/publisher.go`)
- [ ] Email Worker (`internal/workers/email_worker.go`)
- [ ] SMTP / Resend / SendGrid integration + email templates (booking confirmation)

### Phase 5: DevOps & Infrastructure as Code 📋 PLANNED

> Prerequisites: Complete Phase 4

- [ ] Kubernetes: Local cluster (Kind/Minikube) + Helm Charts for GoTicket
- [ ] Terraform: Provision AWS Free Tier resources (RDS, ECS, S3, VPC)
- [ ] Observability: Prometheus metrics + Grafana dashboards + OpenTelemetry distributed tracing

### Phase 6: Production Polish 📋 PLANNED

> Prerequisites: Complete Phase 5

- [ ] Health Checks: `/health`, `/readiness`, `/liveness` endpoints for K8s probes
- [ ] HTTPS/TLS: cert-manager + Let's Encrypt in Kubernetes
- [ ] Secrets Management: K8s Secrets + Sealed Secrets (GitOps-friendly)
- [ ] Error Tracking: Sentry integration (free tier) + error reporting middleware
- [ ] Resilience Patterns:
  - [ ] Circuit Breaker (go-resilience library)
  - [ ] Retry Logic with exponential backoff
  - [ ] Timeout configuration for all external calls

---

## 3. Immediate Next Steps (Priority Order)

1. **Correlation IDs** — essential for debugging distributed systems
2. **RabbitMQ** → **Email Worker** — complete the event-driven pipeline
3. **Kubernetes** (local Kind/Minikube) — deploy the containerized app
5. **Terraform** — infrastructure as code provisioning
6. **Microservices split** — extract BookingService / NotificationService when the domain boundaries are clear

> **Mentorship Rule**: We do NOT copy-paste code. We design the interface first, then the mentee implements it.
