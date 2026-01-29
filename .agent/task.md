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

- [ ] Refactor: DTO Layer (czyste odpowiedzi JSON). <!-- id: 20 -->
- [ ] Refactor: Advanced Error Mapping (Domain Errors). <!-- id: 21 -->
- [ ] Automatyzacja migracji w kodzie (Embedded migrations). <!-- id: 17 -->
- [ ] Testy integracyjne (Database tests). <!-- id: 22 -->

## Phase 5: DevOps & Cloud

- [ ] Konteneryzacja (Dockerfile, Docker Compose dla całości). <!-- id: 13 -->
- [ ] GitHub Actions (CI pipelines). <!-- id: 14 -->
- [ ] Provisioning AWS (Terraform/OpenTofu). <!-- id: 15 -->
- [ ] Deploy na AWS. <!-- id: 16 -->
