# GoTicket - Project Roadmap 🗺️

> **Development roadmap and current status for the GoTicket event ticketing system**

---

## 📊 Current Status

| Phase | Status | Progress |
|-------|--------|----------|
| Phase 1: Core Domain & Repository | ✅ COMPLETE | 100% |
| Phase 2: API Layer & Handlers | ✅ COMPLETE | 100% |
| Phase 3: Authentication & Authorization | 🔄 IN PROGRESS | 90% |
| Phase 4: Security & Infrastructure | 📋 PLANNED | 0% |
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

## Phase 3: Authentication & Authorization 🔄

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

### Remaining
- [ ] Role-based access control (RBAC)
  - [ ] Extract role from JWT claims in middleware
  - [ ] Store role in context (not just user_id)
  - [ ] Create `RequireRole()` middleware
  - [ ] Apply role restrictions:
    - [ ] POST /events → organizer/admin only
    - [ ] PUT /events/{id} → organizer/admin only
    - [ ] DELETE /events/{id} → admin only
- [ ] Helper functions
  - [ ] `GetUserFromContext()` - extract user_id and role
- [ ] Testing
  - [ ] Test role-based access (user vs organizer vs admin)

---

## Phase 4: Security & Infrastructure 📋

> **Prerequisites**: Complete Phase 3

### Rate Limiting
- [ ] Add Redis to docker-compose
- [ ] Implement rate limiting middleware
  - [ ] IP-based: 10 attempts / 15 minutes
  - [ ] Email-based: 5 attempts / 15 minutes
  - [ ] Hybrid approach (both IP and email)
- [ ] Test rate limiting (manual + automated)

### Reverse Proxy (nginx)
- [ ] Configure nginx as reverse proxy
- [ ] SSL/TLS certificates (Let's Encrypt)
- [ ] Request logging and access logs
- [ ] Static file serving (if needed)
- [ ] Load balancing configuration

### Additional Security
- [ ] CORS configuration
- [ ] Request size limits
- [ ] Timeout configuration
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

1. **Role-Based Authorization (RBAC)** - implement role checking middleware
2. **Apply role restrictions** - organizer/admin for event creation, admin for deletion
3. **Testing** - verify role-based access control works correctly
4. **Move to Phase 5** - add unit and integration tests

---

## 📝 Notes

- **Security first**: All auth vulnerabilities (user enumeration, timing attacks) have been addressed
- **Production-ready patterns**: Using bcrypt, JWT, proper error handling
- **Scalability**: Architecture supports Redis and nginx integration
- **Documentation**: All code is documented and follows Go best practices
