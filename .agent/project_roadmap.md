# Mentorship Roadmap: go-ticket Project

This document outlines the evaluation of the current `go-ticket` project and the roadmap to transform it into a Senior-Level Fullstack/DevOps portfolio project.

## 1. Project Audit: Current State vs. Goal

| Category         | Component        | Current State              | Target "Senior" State                            | Status           |
| :--------------- | :--------------- | :------------------------- | :----------------------------------------------- | :--------------- |
| **Architecture** | Pattern          | Modular Monolith (Layered) | Domain-Driven Design (DDD) with clear boundaries | 🟡 In Progress   |
|                  | Separation       | Basic (API, Service, Repo) | Hexagonal / Clean Architecture (Strict)          | 🟢 Good Start    |
| **Backend (Go)** | API              | REST (Standard Lib)        | REST + **GraphQL** + gRPC                        | 🔴 Missing       |
|                  | Concurrency      | Basic Mutexes              | Advanced Patterns (Workers, Pipelines)           | 🟡 Basic         |
|                  | Persistence      | Postgres (pgx)             | Postgres + **Redis** (Caching)                   | 🟡 Postgres only |
| **Security**     | Authentication   | None                       | **JWT** + **OAuth2** (Google/GitHub)             | 🔴 Missing       |
|                  | Authorization    | None                       | **RBAC** (Role-Based Access Control)             | 🔴 Missing       |
| **DevOps**       | Containerization | `docker-compose` (Dev)     | Optimized Multi-stage **Dockerfiles**            | 🔴 Missing       |
|                  | Orchestration    | None                       | **Kubernetes** (Helm/Kustomize)                  | 🔴 Missing       |
|                  | IaC              | None                       | **Terraform** / OpenTofu                         | 🔴 Missing       |
|                  | CI/CD            | None                       | **GitHub Actions** (Lint, Test, Build, Push)     | 🔴 Missing       |
| **Messaging**    | Async            | None                       | **Kafka / RabbitMQ** (Event Driven)              | 🔴 Missing       |
| **Quality**      | Testing          | Unit Tests (Domain only)   | **Integration**, **E2E**, Property-based Tests   | 🔴 Critical Gap  |
|                  | Observability    | Logging (slog)             | Distributed Tracing (OTEL), Metrics (Prometheus) | 🟡 Logging only  |

## 2. The "Antigravity" Roadmap

We will not build everything at once. We will follow an iterative "Evolutionary Architecture" approach.

### Phase 1: Foundation Hardening (The "Professional" Monolith)

Before adding complexity, we must ensure quality.

- [ ] **Testing Strategy**: Add Integration Tests for Repositories (Testcontainers) and API Handlers.
- [ ] **Observability**: Add structured logging with correlation IDs and basic metrics.
- [ ] **Configuration**: Robust config management (Viper or strict env parsing).
- [ ] **Documentation**: Generate **Swagger/OpenAPI** docs (using `swaggo` or similar).
- [ ] **Automation**: Create a `Makefile` for common tasks (build, test, lint, docker-up).

### Phase 2: Advanced Backend Features

Expand the application capabilities to learn modern API standards.

- [ ] **GraphQL**: Implement a GraphQL layer using `gqlgen`.
- [ ] **Caching**: Introduce Redis for caching event details.
- [ ] **Security & Auth**:
  - [ ] User Registration & Login endpoints.
  - [ ] JWT-based authentication (access + refresh tokens).
  - [ ] OAuth2 integration (Google, GitHub).
  - [ ] RBAC middleware (admin/organizer/user roles).
  - [ ] Password security (bcrypt, cost 12+).
- [ ] **Security & Observability**:
  - [ ] Rate Limiting Middleware (protect against scanning attacks).
  - [ ] Audit Logging Middleware (compliance: GDPR, PCI DSS).
  - [ ] Metrics collection (request count, error rate, latency).
  - [ ] Correlation IDs for request tracing.

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

1.  **Add Integration Tests**: The project has only one domain test. This is the biggest risk. We need to test the Repository layer with a real DB (using `testcontainers-go` is the senior way).
2.  **Dockerize app**: Create a `Dockerfile` to run the app in isolation, not just the DB.
3.  **GraphQL Definition**: Start defining the schema for the query side.

> [!IMPORTANT]
> **Mentorship Rule**: We will NOT copy-paste code. We will design the interface first, then you implement it.
