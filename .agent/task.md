# Mentoring Plan: Project "GoTicket"

## Globalne Cele

- [ ] Opanowanie zaawansowanego Go (Concurrency, Interfaces, Reflect).
- [ ] Zrozumienie architektury mikroserwisowej i komunikacji gRPC.
- [ ] Praktyczne zastosowanie wzorców projektowych i DDD/Hexagonal.
- [ ] DevOps: Docker, CI/CD, AWS, Monitoring.

## Phase 0: Setup & Foundation

- [x] Wybór i zatwierdzenie tematu projektu. <!-- id: 0 -->
- [x] Konfiguracja Repozytorium (Monorepo vs Polyrepo - decyzja: Monorepo). <!-- id: 1 -->
- [x] Setup środowiska lokalnego (Go, Docker, Makefiles). <!-- id: 2 -->
- [x] Projekt architektury wysokopoziomowej (C4 Model). <!-- id: 3 -->

## Phase 1: The Core (Prosty Monolit modularny na start -> refactor do Micro)

_Uwaga: Zaczniemy od modularnego monolitu, żeby skupić się na kodzie domenowym, a potem rozerwiemy go na mikroserwisy. To najzdrowsza ścieżka rozwoju._

- [x] Struktura projektu (Go Standard Layout). <!-- id: 4 -->
- [x] Implementacja domeny `Events` (Structs, Entities). <!-- id: 5 -->
- [x] Implementacja warstwy dostępu do danych (Postgres + pgx/sqlc). <!-- id: 6 -->

## Phase 2: API & Communication

- [x] Wystawienie REST API (net/http). <!-- id: 7 -->
- [ ] Middleware (Logging, Recovery, Auth). <!-- id: 8 -->
- [ ] Wprowadzenie gRPC (Proto definitions). <!-- id: 9 -->

## Phase 3: The Hard Parts (Concurrency)

- [ ] System rezerwacji (Booking Logic). <!-- id: 10 -->
- [ ] Implementacja blokowania zasobów (Mutex vs DB Locks). <!-- id: 11 -->
- [ ] Testy obciążeniowe (k6) i optymalizacja. <!-- id: 12 -->

## Phase 4: DevOps & Cloud

- [ ] Automatyzacja migracji w kodzie (Embedded migrations). <!-- id: 17 -->
- [ ] Konteneryzacja (Dockerfile, Docker Compose dla całości). <!-- id: 13 -->
- [ ] GitHub Actions (CI pipelines). <!-- id: 14 -->
- [ ] Provisioning AWS (Terraform/OpenTofu). <!-- id: 15 -->
- [ ] Deploy na AWS. <!-- id: 16 -->
