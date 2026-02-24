# GoTicket - Project Roadmap 🗺️

> **Development roadmap and current status for the GoTicket event ticketing system**

---

## 📊 Current Status

| Phase | Status | Progress |
|-------|--------|----------|
| Phase 1: Core Domain & Repository | ✅ COMPLETE | 100% |
| Phase 2: API Layer & Handlers | ✅ COMPLETE | 100% |
| Phase 3: Authentication & Authorization | ✅ COMPLETE | 100% |
| Phase 4: Security & Infrastructure | 🔄 IN PROGRESS | 30% |
| Phase 5: Testing & CI/CD | 📋 PLANNED | 0% |
| Phase 6: Production Deployment | 📋 PLANNED | 0% |

---

## Phase 1: Core Domain & Repository ✅

### Domain Models
- [x] Event (with business logic: capacity, booking validation)
- [x] Booking (with status management)
- [x] User (with password hashing, role management)
- [x] Domain errors (sentinel errors for all entities)

### Repository Layer
- [x] EventRepository (CRUD operations)
- [x] BookingRepository (CRUD operations)
- [x] UserRepository (CRUD + GetByEmail, GetByID)
- [x] sqlc integration (type-safe SQL queries)
- [x] PostgreSQL migrations

---

## Phase 2: API Layer & Handlers ✅

### HTTP Handlers
- [x] Event handlers (Create, Update, Delete, Get, List)
- [x] Booking handlers (Create)
- [x] Response helpers (error handling, JSON responses)
- [x] Error mapping (domain errors → HTTP status codes)

### DTOs
- [x] Event DTOs (request/response)
- [x] Booking DTOs (request/response)
- [x] Auth DTOs (register/login)

---

## Phase 3: Authentication & Authorization ✅

### Completed ✅
- [x] JWT service (token generation/verification)
- [x] Password hashing (bcrypt with proper cost)
- [x] User registration endpoint (`POST /auth/register`)
- [x] User login endpoint (`POST /auth/login`)
- [x] Security fixes:
  - [x] User enumeration protection
  - [x] Timing attack protection
  - [x] Password validation (min 8 characters)
- [x] Response helpers (DRY - exported functions)
- [x] Error mapping for auth errors
- [x] **Auth middleware (JWT verification)**
- [x] **Protected endpoints (all event/booking endpoints)**
- [x] **Security logging (IP, path, no token exposure)**
- [x] **RBAC — `RequireRole()` middleware**
- [x] **Role restrictions applied:**
  - [x] `POST /events` → organizer only
  - [x] `PUT /events/{id}` → organizer only
  - [x] `DELETE /events/{id}` → admin only
  - [x] `GET /events`, `GET /events/{id}`, `POST bookings` → all authenticated users

---

## Phase 4: Security & Infrastructure 🔄

> **Prerequisites**: Complete Phase 3

### Containerization
- [x] Multi-stage Dockerfile (golang:alpine builder → distroless/static-debian12 final)
- [x] Docker Compose: app + postgres services, healthcheck, depends_on, env_file
- [x] Image size: 31MB (vs ~600MB builder)
- [x] Non-root user (nonroot:nonroot) in container

### CI/CD
- [/] GitHub Actions workflow (lint → test → build → push) ← IN PROGRESS

### Rate Limiting
- [ ] Add Redis to docker-compose
- [ ] Implement rate limiting middleware
  - [ ] IP-based: 10 attempts / 15 minutes
  - [ ] Email-based: 5 attempts / 15 minutes

### Additional Security
- [ ] CORS configuration
- [ ] Request size limits
- [ ] Security headers (CSP, HSTS, etc.)

---

## Phase 5: Testing & CI/CD 📋

> **Prerequisites**: Complete Phase 3 & 4

### Testing
- [ ] Unit tests (domain, service layers)
- [ ] Integration tests (repository, API)
- [ ] E2E tests (full user flows)
- [ ] Test coverage > 80%

### CI/CD
- [ ] GitHub Actions workflow
- [ ] Automated testing on PR
- [ ] Linting and formatting checks
- [ ] Docker image build
- [ ] Deployment automation

---

## Phase 6: Production Deployment 📋

> **Prerequisites**: Complete Phase 5

### Infrastructure
- [ ] Cloud provider setup (AWS/GCP/DigitalOcean)
- [ ] Database hosting (managed PostgreSQL)
- [ ] Redis hosting (managed Redis)
- [ ] Container orchestration (Docker Swarm/Kubernetes)

### Monitoring & Observability
- [ ] Prometheus metrics
- [ ] Grafana dashboards
- [ ] Alerting rules
- [ ] Log aggregation (Loki/ELK)

### Production Readiness
- [ ] Health checks (`/health`, `/ready`)
- [ ] Graceful shutdown
- [ ] Database connection pooling
- [ ] Backup and recovery procedures

---

## 🎯 Next Steps (Immediate)

1. **Rate Limiting** - Redis-based, IP+email hybrid (Phase 4)
2. **CORS & Security Headers** - CSP, HSTS, request size limits
3. **Load Testing (k6)** - verify race conditions under load
4. **Move to Phase 5** - GitHub Actions CI/CD pipeline

---

## 📝 Notes

- **Security first**: All auth vulnerabilities (user enumeration, timing attacks) have been addressed
- **Production-ready patterns**: Using bcrypt, JWT, proper error handling
- **Scalability**: Architecture supports Redis and nginx integration
- **Documentation**: All code is documented and follows Go best practices
