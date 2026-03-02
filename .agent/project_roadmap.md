# Mentorship Roadmap: go-ticket Project

This document outlines the evaluation of the current `go-ticket` project and the roadmap to transform it into a Senior-Level Fullstack/DevOps portfolio project.

## 1. Project Audit: Current State vs. Goal

| Category         | Component        | Current State                         | Target "Senior" State                            | Status             |
| :--------------- | :--------------- | :------------------------------------ | :----------------------------------------------- | :----------------- |
| **Architecture** | Pattern          | Modular Monolith (Clean Architecture) | Domain-Driven Design (DDD) with clear boundaries | 🟢 Solid           |
|                  | Separation       | Domain / Service / Repo / API         | Hexagonal / Clean Architecture (Strict)          | 🟢 Done            |
| **Backend (Go)** | API              | REST (Standard Lib)                   | REST + **GraphQL** + gRPC                        | 🟡 REST done       |
|                  | Concurrency      | Atomic booking, Testcontainers        | Advanced Patterns (Workers, Pipelines)           | 🟢 Good            |
|                  | Persistence      | Postgres (pgx + sqlc)                 | Postgres + **Redis** (Caching + Rate Limiting)   | 🟡 Postgres done   |
| **Security**     | Authentication   | JWT (access token)                    | **JWT** + **OAuth2** (Google/GitHub)             | 🟢 JWT done        |
|                  | Authorization    | RBAC (user/organizer/admin)           | Full **RBAC** + fine-grained permissions         | 🟢 Done            |
| **DevOps**       | Containerization | Multi-stage Dockerfile (distroless)   | Optimized Multi-stage **Dockerfiles**            | 🟢 Done            |
|                  | Orchestration    | None                                  | **Kubernetes** (Helm/Kustomize)                  | 🔴 Planned         |
|                  | IaC              | None                                  | **Terraform** / OpenTofu                         | 🔴 Planned         |
|                  | CI/CD            | GitHub Actions (CI + CD)              | **GitHub Actions** (Lint, Test, Build, Push)     | 🟢 Done            |
| **Messaging**    | Async            | None                                  | **Kafka / RabbitMQ** (Event Driven)              | 🔴 Phase 5         |
| **Quality**      | Testing          | Integration Tests (Testcontainers)    | **E2E**, Load (k6), Property-based Tests         | 🟢 Int. tests done |
|                  | Observability    | Structured logging (slog/JSON)        | Distributed Tracing (OTEL), Metrics (Prometheus) | 🟡 Logging only    |

## 2. The "Antigravity" Roadmap

We will not build everything at once. We will follow an iterative "Evolutionary Architecture" approach.

### Phase 1: Foundation Hardening ✅ COMPLETE

- [x] **Testing Strategy**: Integration Tests (Testcontainers) for Repositories + API Handlers.
- [x] **Structured Logging**: JSON logging via `slog`.
- [x] **Security & Auth**: JWT, RBAC (user/organizer/admin), bcrypt, user enumeration protection.
- [x] **Error Handling**: Domain errors → HTTP status codes mapping.
- [x] **Embedded Migrations**: Auto-run on startup.

### Phase 2: Security & DevOps Foundations 🔄 CURRENT

- [x] **Dockerfile**: Multi-stage build (builder → distroless/static-debian12). 31MB image ✅
- [x] **docker-compose**: app + postgres services, healthcheck, env_file. ✅
- [x] **GitHub Actions CI**: lint (golangci-lint v2) + test (-race, coverage) + build on every push. ✅
- [x] **GitHub Actions CD**: docker build + push to GHCR on main (workflow_run). ✅
- [x] **Code Quality**: Fixed 36 golangci-lint issues (errcheck, gosec, lll, gocyclo, funlen). ✅
- [x] **Rate Limiting**: Redis INCR+EXPIRE, two limiters (auth: 5/15min IP-based, API: 100/1min user-based), miniredis unit tests. ✅
- [x] **main.go refactor**: extracted setupRoutes, setupServer, gracefulShutdown, setupRepositories helpers. ✅
- [x] **MapDomainError refactor**: table-driven lookup with errorMapping struct. ✅
- [/] **Load Testing (k6)**: Verify race-condition safety under load.
  - ✅ `tests/load/seed.sql` — organizer user + event z capacity=10000
  - ✅ `tests/load/booking_scenario.js` — szkielet: options, stages, thresholds, default()
  - ✅ `tests/load/README.md` — dokumentacja uruchomienia, seed, weryfikacja, CI/CD plan
  - ⏳ Dokończyć `booking_scenario.js`: `setup()` (login→token) + `default(data)` z auth headerem
  - ⏳ Uruchomić seed + test + zweryfikować brak double-bookingu
- [ ] **Swagger/OpenAPI**: Auto-generated docs (`swaggo`).
- [ ] **Correlation IDs**: Request tracing through middleware.

### Phase 3: DevOps & Containerization

Move from "running locally" to "production ready".

- [ ] **Docker**: Create optimized, multi-stage Dockerfiles (Distroless).
- [ ] **CI Pipeline**: GitHub Actions for automatic linting and testing.
- [ ] **Registry**: Pushing images to a container registry.

### Phase 4: Infrastructure as Code (IaC) & Orchestration

The "DevOps" heavy lifting.

- [ ] **Terraform**: Provision local infra (Docker based) or AWS Free Tier resources (`RDS`, `S3`).
- [ ] **Kubernetes**: Deploy the app to a local Cluster (Kind/Minikube) using Manifests and Helm.

### Phase 5: Microservices & Event-Driven Architecture

Refactoring the monolith into distributed services.

- [ ] **Split**: Extract `BookingService` or `NotificationService`.
- [ ] **Async Messaging**: Implement RabbitMQ/Kafka for "BookingConfirmed" events.
- [ ] **Saga Pattern**: Handle distributed transactions (if needed).

### Phase 6: Production Polish (High-Impact Additions)

Add critical production features that demonstrate production-ready mindset.

- [ ] **Health Checks**: `/health`, `/readiness`, `/liveness` endpoints for K8s probes.
- [ ] **HTTPS/TLS**: cert-manager + Let's Encrypt integration in K8s.
- [ ] **Secrets Management**: K8s Secrets + Sealed Secrets (GitOps-friendly).
- [ ] **Error Tracking**: Sentry integration (free tier) for runtime error monitoring.
- [ ] **Resilience Patterns**:
  - [ ] Circuit Breaker (go-resilience library).
  - [ ] Retry Logic with exponential backoff.
  - [ ] Timeout configuration for external calls.

> [!NOTE]
> **Why these 6 features?**
>
> - **High Impact**: Every production system has these
> - **Low Effort**: ~15-25 hours total (2-3 weeks part-time)
> - **Learning Value**: K8s probes, TLS, secrets, error tracking, resilience
> - **CV/Portfolio**: Shows production mindset, not just "tutorial code"

## 3. Immediate Next Steps / Recommendations

1. **Dockerfile** (multi-stage, distroless) — fundament wszystkiego co dalej (K8s, CI/CD). **← TERAZ**
2. **Rate Limiting** (Redis) — security-first + nauka Redis przy okazji.
3. **GitHub Actions** — CI/CD pipeline (lint → test → build → push image).
4. **Kubernetes** (local Kind/Minikube) — deploy skonteneryzowanej app.
5. **Terraform** — provisioning infrastruktury.
6. **Microservices split** — wyciągnięcie BookingService / NotificationService.
7. **Kafka / RabbitMQ** — event-driven communication między serwisami.

> [!IMPORTANT]
> **Mentorship Rule**: We will NOT copy-paste code. We design the interface first, then you implement it.
