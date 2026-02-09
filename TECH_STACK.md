# GoTicket - Tech Stack & Concepts Dictionary 📚

> **Complete guide to technologies and patterns used in the GoTicket project**
>
> This document explains **what**, **why**, and **how** we use each technology in the GoTicket project.

---

## 📑 Table of Contents

1. [Core Technologies](#core-technologies)
2. [Architecture Patterns](#architecture-patterns)
3. [Database Layer](#database-layer)
4. [Testing Infrastructure](#testing-infrastructure)
5. [DevOps & Tooling](#devops--tooling)
6. [Production Patterns](#production-patterns)

---

## Core Technologies

### Go (Golang) 1.22+

**What is it?**

- Programming language created by Google
- Compiled (not interpreted like Python)
- Statically typed
- Built-in concurrency (goroutines)

**Why Go?**

- ✅ **Performance** - Compiles to native code (fast as C++)
- ✅ **Concurrency** - Goroutines handle thousands of requests simultaneously
- ✅ **Simplicity** - Small language (25 keywords), easy to learn
- ✅ **Standard Library** - HTTP server, JSON, crypto built-in
- ✅ **Single Binary** - Deployment is one file (no dependencies)

**Real World:**

- Used by: Google, Uber, Netflix, Dropbox, Docker, Kubernetes
- Ideal for: Microservices, APIs, CLI tools, Cloud infrastructure

**In our project:**

```go
// Example: Goroutines for concurrent bookings
for i := 0; i < 200; i++ {
    go func() {
        bookingService.CreateBooking(...)  // 200 concurrent requests!
    }()
}
```

---

### PostgreSQL 16

**What is it?**

- Relational database (SQL)
- Open-source
- ACID compliant (transactions!)

**Why Postgres?**

- ✅ **ACID Transactions** - Atomicity, Consistency, Isolation, Durability
- ✅ **Advanced Features** - JSON columns, full-text search, GIS
- ✅ **Performance** - Faster than MySQL for complex queries
- ✅ **Reliability** - Production-proven (used by Instagram, Spotify)
- ✅ **Row-Level Locking** - Concurrent updates without deadlocks

**Real World:**

- Used by: Instagram (billions of rows), Spotify, Reddit, Twitch
- Ideal for: Transactional systems, Analytics, Geospatial data

**In our project:**

```sql
-- Atomic spot reservation (prevents double-booking)
UPDATE events
SET available_spots = available_spots - $1
WHERE id = $2 AND available_spots >= $1;
-- ↑ If insufficient spots, UPDATE returns 0 rows (transaction rollback!)
```

---

### `net/http` (Standard Library)

**What is it?**

- HTTP server built into Go
- Zero external dependencies
- Production-ready

**Why Standard Library instead of Gin/Fiber?**

- ✅ **Zero Magic** - Understand how HTTP works under the hood
- ✅ **Stability** - Doesn't change (backward compatible)
- ✅ **Performance** - No framework overhead
- ✅ **Go 1.22 Routing** - New features (path parameters, HTTP methods)

**Real World:**

- Used by: Cloudflare, Twitch API, HashiCorp tools
- Ideal for: APIs, Microservices, when you need control

**In our project:**

```go
// Go 1.22 routing - elegant and readable
mux.HandleFunc("POST /events", handler.CreateEvent)
mux.HandleFunc("GET /events/{id}", handler.GetEvent)
mux.HandleFunc("POST /events/{event_id}/bookings", handler.CreateBooking)
```

---

## Architecture Patterns

### Clean Architecture (Hexagonal Architecture)

**What is it?**

- Layer separation: Domain → Repository → Service → API
- Domain logic independent of infrastructure
- Dependency Inversion Principle

**Why Clean Architecture?**

- ✅ **Testability** - Domain logic without database
- ✅ **Flexibility** - Change DB/API without changing logic
- ✅ **Maintainability** - Each layer has clear responsibility

**Layers in project:**

```
┌─────────────────────────────────────┐
│   API Layer (internal/api)          │  ← HTTP handlers, DTOs, JSON
├─────────────────────────────────────┤
│   Service Layer (internal/services) │  ← Transactions, orchestration
├─────────────────────────────────────┤
│   Repository (internal/postgres)    │  ← Database queries (sqlc)
├─────────────────────────────────────┤
│   Domain Layer (internal/domain)    │  ← Business logic, validation
└─────────────────────────────────────┘
```

**Real World:**

- Used by: Netflix, Uber, Spotify backend teams
- Ideal for: Long-lived projects, teams >5 people

**In our project:**

```go
// Domain - pure business logic (no DB, no HTTP)
type Event struct {
    id             uuid.UUID
    availableSpots int
}

func (e *Event) ReserveSpots(count int) error {
    if e.availableSpots < count {
        return ErrEventIsFull  // Domain error!
    }
    e.availableSpots -= count
    return nil
}

// Repository - only DB operations
func (r *EventRepository) ReserveSpots(ctx, eventID, spots) error {
    // SQL query
}

// Service - orchestration + transactions
func (s *BookingService) CreateBooking(ctx, ...) error {
    tx.Begin()
    defer tx.Rollback()

    // Use repository
    err := s.eventRepo.ReserveSpots(...)
    err = s.bookingRepo.CreateBooking(...)

    tx.Commit()
}

// API - HTTP handling
func (h *HTTPHandler) CreateBooking(w, r) {
    // Parse JSON
    // Call service
    // Return JSON response
}
```

---

### Domain-Driven Design (DDD)

**What is it?**

- Business logic at the center of architecture
- Entities, Value Objects, Aggregates
- Ubiquitous Language (business language in code)

**Why DDD?**

- ✅ **Business Focus** - Code reflects business reality
- ✅ **Encapsulation** - Data + logic together (not anemic models)
- ✅ **Validation** - Impossible to create invalid entity

**Real World:**

- Used by: Amazon, Microsoft Azure teams
- Ideal for: Complex business domains, e-commerce, booking systems

**In our project:**

```go
// ❌ Anemic Model (BAD)
type Event struct {
    ID             uuid.UUID
    AvailableSpots int  // Public! Anyone can change!
}
event.AvailableSpots = -100  // Invalid state!

// ✅ Rich Domain Model (GOOD)
type Event struct {
    id             uuid.UUID  // Private!
    availableSpots int        // Private!
}

func NewEvent(name string, capacity int) (*Event, error) {
    if capacity < 0 {
        return nil, errors.New("capacity must be positive")
    }
    return &Event{
        id:             uuid.New(),
        availableSpots: capacity,
    }, nil
}

func (e *Event) AvailableSpots() int {
    return e.availableSpots  // Read-only getter
}

func (e *Event) ReserveSpots(count int) error {
    if count <= 0 {
        return errors.New("count must be positive")
    }
    if e.availableSpots < count {
        return ErrEventIsFull
    }
    e.availableSpots -= count  // Controlled mutation!
    return nil
}
```

**Benefits:**

- Impossible to create `Event` with capacity < 0
- Impossible to reserve -5 spots
- Business logic in one place (not scattered)

---

### Repository Pattern

**What is it?**

- Abstraction over database
- Interface defines operations (GetEvent, CreateBooking)
- Implementation uses specific database (Postgres, MySQL, MongoDB)

**Why Repository Pattern?**

- ✅ **Testability** - Mock repository in tests
- ✅ **Flexibility** - Change DB without changing service layer
- ✅ **Separation** - Domain doesn't know about SQL

**Real World:**

- Used by: All large projects (standard pattern)
- Ideal for: Every project with database

**In our project:**

```go
// Interface (contract)
type EventRepository interface {
    GetEvent(ctx context.Context, id uuid.UUID) (*domain.Event, error)
    CreateEvent(ctx context.Context, event *domain.Event) error
    ReserveSpots(ctx context.Context, eventID uuid.UUID, spots int) error
}

// Postgres implementation
type PostgresEventRepository struct {
    queries *Queries
}

func (r *PostgresEventRepository) GetEvent(ctx, id) (*domain.Event, error) {
    row, err := r.queries.GetEvent(ctx, id)
    // Convert DB row → domain.Event
}

// Service uses interface (not concrete implementation!)
type BookingService struct {
    eventRepo EventRepository  // ← Interface, not *PostgresEventRepository!
}

// In tests: mock repository
type MockEventRepository struct {}
func (m *MockEventRepository) GetEvent(...) (*domain.Event, error) {
    return &domain.Event{...}, nil  // Fake data, no DB!
}
```

---

## Database Layer

### `pgx/v5` + `pgxpool`

**What is it?**

- `pgx` - Pure Go Postgres driver
- `pgxpool` - Connection pooling

**Why pgx instead of database/sql?**

- ✅ **Performance** - 2-3x faster than `database/sql`
- ✅ **Features** - LISTEN/NOTIFY, COPY, binary protocol
- ✅ **Type Safety** - Native Postgres types (UUID, JSONB, arrays)
- ✅ **Connection Pooling** - Built-in, production-ready

**Real World:**

- Used by: Grafana, CockroachDB, Supabase
- Ideal for: High-performance APIs, real-time systems

**In our project:**

```go
// Connection pool (reuse connections)
pool, err := pgxpool.New(ctx, databaseURL)

// Pool automatically:
// - Reuses connections (doesn't create new one each time)
// - Handles reconnects (if DB down)
// - Limits max connections (doesn't overload DB)

// Query
row, err := pool.QueryRow(ctx, "SELECT * FROM events WHERE id = $1", eventID)
```

**Connection Pool - why?**

```
❌ Without pool:
Request 1 → New DB connection (slow!)
Request 2 → New DB connection (slow!)
Request 3 → New DB connection (slow!)

✅ With pool:
Request 1 → Reuse connection #1 (fast!)
Request 2 → Reuse connection #2 (fast!)
Request 3 → Reuse connection #1 (fast!)
```

---

### `sqlc` - SQL Code Generator

**What is it?**

- Generates type-safe Go code from SQL queries
- Compile-time safety (SQL errors during compilation!)
- Zero reflection, zero ORM magic

**Why sqlc instead of GORM/Ent?**

- ✅ **Performance** - No ORM overhead (raw SQL)
- ✅ **Control** - Write exactly the SQL you want
- ✅ **Type Safety** - SQL errors during `go generate`, not runtime!
- ✅ **Simplicity** - Understand what happens (no magic)

**Real World:**

- Used by: Grafana, PlanetScale, Railway
- Ideal for: Performance-critical apps, when you know SQL

**In our project:**

**1. Write SQL (`queries/events.sql`):**

```sql
-- name: GetEvent :one
SELECT * FROM events WHERE id = $1;

-- name: ReserveSpots :exec
UPDATE events
SET available_spots = available_spots - $1
WHERE id = $2 AND available_spots >= $1;
```

**2. sqlc generates Go code:**

```go
// Auto-generated!
type GetEventRow struct {
    ID             pgtype.UUID
    Name           string
    AvailableSpots int32
}

func (q *Queries) GetEvent(ctx context.Context, id pgtype.UUID) (GetEventRow, error) {
    // SQL query implementation
}

func (q *Queries) ReserveSpots(ctx context.Context, arg ReserveSpotsParams) error {
    // SQL query implementation
}
```

**3. Use in repository:**

```go
row, err := r.queries.GetEvent(ctx, eventID)
// ↑ Type-safe! Compiler knows GetEvent returns GetEventRow
```

**Benefits:**

- Typo in SQL? → Error during `sqlc generate`
- Schema change? → Compilation error (not runtime!)
- Performance? → Raw SQL (no ORM overhead)

---

### `golang-migrate` - Database Migrations

**What is it?**

- CLI tool for managing migrations
- Versioned schema changes
- Up/Down migrations

**Why migrations?**

- ✅ **Version Control** - Schema in Git (like code!)
- ✅ **Reproducibility** - Every developer has same schema
- ✅ **Rollback** - Down migrations revert changes
- ✅ **Production Safety** - Atomic changes, tracking table

**Real World:**

- Standard in every project with database
- Used by: All production apps

**In our project:**

**Migration files:**

```
migrations/
├── 000001_create_events_table.up.sql    ← Apply
├── 000001_create_events_table.down.sql  ← Rollback
├── 000002_add_bookings.up.sql
└── 000002_add_bookings.down.sql
```

**Up migration (apply):**

```sql
-- 000001_create_events_table.up.sql
CREATE TABLE events (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    capacity INT NOT NULL
);
```

**Down migration (rollback):**

```sql
-- 000001_create_events_table.down.sql
DROP TABLE events;
```

**Tracking table (`schema_migrations`):**

```sql
SELECT * FROM schema_migrations;
-- version | dirty
-- --------+-------
--    2    | f      ← Currently at version 2, no errors
```

**Embedded migrations (`//go:embed`):**

```go
//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(databaseURL string) error {
    // Migrations embedded in binary!
    // No external files needed in production
}
```

---

## Testing Infrastructure

### Testcontainers

**What is it?**

- Library for running Docker containers in tests
- Real Postgres instance for each test
- Automatic cleanup

**Why Testcontainers instead of mocks?**

- ✅ **Real Database** - Test real SQL, not fake behavior
- ✅ **Integration Tests** - Catch bugs mocks miss
- ✅ **Isolation** - Each test has clean database
- ✅ **CI/CD Ready** - Works in GitHub Actions, GitLab CI

**Real World:**

- Used by: Spring Boot, Quarkus, .NET teams
- Ideal for: Integration tests, E2E tests

**In our project:**

```go
func TestEventRepository_CreateEvent(t *testing.T) {
    ctx := context.Background()

    // Start real Postgres container
    container, _ := postgres.Run(ctx, "postgres:16", ...)
    defer container.Terminate(ctx)  // Auto-cleanup!

    // Get connection string
    connStr, _ := container.ConnectionString(ctx)

    // Connect to real database
    pool, _ := pgxpool.New(ctx, connStr)

    // Test against REAL Postgres!
    repo := NewEventRepository(pool)
    err := repo.CreateEvent(ctx, event)

    // Verify in database
    row := pool.QueryRow("SELECT * FROM events WHERE id = $1", event.ID())
}
```

**Benefits:**

- Catch SQL syntax errors (mocks won't catch!)
- Test transactions, locks, constraints
- Confidence that code works in production

---

### Table-Driven Tests

**What is it?**

- Testing pattern in Go
- One test function, multiple test cases
- Readable and maintainable

**Why Table-Driven Tests?**

- ✅ **DRY** - Don't duplicate test code
- ✅ **Readability** - Test cases in table (easy to add new)
- ✅ **Coverage** - Easy to cover edge cases

**Real World:**

- Standard in Go community
- Used by: Go standard library, all Go projects

**In our project:**

```go
func TestEvent_ReserveSpots(t *testing.T) {
    tests := []struct {
        name          string
        capacity      int
        reserveCount  int
        wantErr       error
        wantRemaining int
    }{
        {
            name:          "success - reserve 10 from 100",
            capacity:      100,
            reserveCount:  10,
            wantErr:       nil,
            wantRemaining: 90,
        },
        {
            name:          "error - insufficient spots",
            capacity:      5,
            reserveCount:  10,
            wantErr:       domain.ErrEventIsFull,
            wantRemaining: 5,
        },
        {
            name:          "error - negative count",
            capacity:      100,
            reserveCount:  -5,
            wantErr:       domain.ErrInvalidCount,
            wantRemaining: 100,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            event, _ := domain.NewEvent("Test", tt.capacity)
            err := event.ReserveSpots(tt.reserveCount)

            assert.Equal(t, tt.wantErr, err)
            assert.Equal(t, tt.wantRemaining, event.AvailableSpots())
        })
    }
}
```

**Output:**

```
=== RUN   TestEvent_ReserveSpots
=== RUN   TestEvent_ReserveSpots/success_-_reserve_10_from_100
=== RUN   TestEvent_ReserveSpots/error_-_insufficient_spots
=== RUN   TestEvent_ReserveSpots/error_-_negative_count
--- PASS: TestEvent_ReserveSpots (0.00s)
```

---

### Race Condition Testing

**What is it?**

- Test concurrent access (goroutines)
- Detect data races, deadlocks
- `go test -race` flag

**Why Race Tests?**

- ✅ **Concurrency Bugs** - Catch double-booking, lost updates
- ✅ **Production Confidence** - Know code works under load
- ✅ **Database Locks** - Verify row-level locking works

**Real World:**

- Critical for: Booking systems, payment processing, inventory
- Used by: Uber, Airbnb, Ticketmaster

**In our project:**

```go
func TestEventRepository_ReserveSpots_Concurrent(t *testing.T) {
    event := createTestEvent(ctx, t, pool, WithCapacity(100))

    // 200 goroutines try to reserve 1 spot each
    numGoroutines := 200
    var wg sync.WaitGroup
    successCount := atomic.Int32{}

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            err := repo.ReserveSpots(ctx, event.ID(), 1)
            if err == nil {
                successCount.Add(1)  // Thread-safe increment
            }
        }()
    }

    wg.Wait()

    // Only 100 should succeed (not 200!)
    assert.Equal(t, int32(100), successCount.Load())

    // Verify database state
    retrieved := getEventFromDB(ctx, t, pool, event.ID())
    assert.Equal(t, 0, retrieved.AvailableSpots())
}
```

**What we test:**

- 200 goroutines → 100 spots → only 100 succeed
- Database locks prevent double-booking
- `atomic.Int32` for thread-safe counting

**Run with race detector:**

```bash
go test -race ./internal/postgres/...
# Detects race conditions (if any!)
```

---

## DevOps & Tooling

### Docker & Docker Compose

**What is it?**

- Docker - Containerization platform
- Docker Compose - Multi-container orchestration

**Why Docker?**

- ✅ **Consistency** - "Works on my machine" → Works everywhere
- ✅ **Isolation** - Postgres in container (doesn't pollute system)
- ✅ **Reproducibility** - Every developer has same env
- ✅ **Production Parity** - Dev environment = Production

**Real World:**

- Industry standard (every company uses it)
- Ideal for: Development, CI/CD, Production deployment

**In our project:**

**`docker-compose.yml`:**

```yaml
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: go_ticket
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

**Usage:**

```bash
# Start Postgres
docker compose up -d

# Stop Postgres
docker compose down

# Remove data (fresh start)
docker compose down -v
```

---

### `golangci-lint`

**What is it?**

- Meta-linter (runs 50+ linters)
- Static code analysis
- Catches bugs, style issues, performance problems

**Why linter?**

- ✅ **Code Quality** - Enforce best practices
- ✅ **Bug Detection** - Catch errors before runtime
- ✅ **Consistency** - Team follows same style
- ✅ **Learning** - Linter teaches better patterns

**Real World:**

- Standard in every Go project
- Used by: Google, Uber, all open-source projects

**In our project:**

**`.golangci.yml`:**

```yaml
linters:
  enable:
    - errcheck # Check error handling
    - gosimple # Simplify code
    - govet # Go vet checks
    - ineffassign # Detect useless assignments
    - staticcheck # Advanced static analysis
    - unused # Detect unused code
```

**Run:**

```bash
golangci-lint run
# ✅ All checks passed!
```

**Example catches:**

```go
// ❌ Linter error: error not checked
pool.QueryRow(ctx, "SELECT ...")

// ✅ Fixed
row, err := pool.QueryRow(ctx, "SELECT ...")
if err != nil {
    return err
}
```

---

### `slog` - Structured Logging

**What is it?**

- Standard library logger (Go 1.21+)
- Structured logging (JSON format)
- Levels: DEBUG, INFO, WARN, ERROR

**Why slog instead of fmt.Println?**

- ✅ **Structured** - JSON output (parseable by tools)
- ✅ **Levels** - Filter logs in production
- ✅ **Context** - Add metadata (user_id, request_id)
- ✅ **Performance** - Zero-allocation logging

**Real World:**

- Standard for production apps
- Ideal for: Observability, debugging, monitoring

**In our project:**

```go
// Setup
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
slog.SetDefault(logger)

// Usage
slog.Info("Server started", "port", 8080)
slog.Error("Database error", "error", err, "query", "SELECT ...")

// Output (JSON):
// {"time":"2024-01-01T12:00:00Z","level":"INFO","msg":"Server started","port":8080}
// {"time":"2024-01-01T12:00:01Z","level":"ERROR","msg":"Database error","error":"connection refused","query":"SELECT ..."}
```

**Benefits:**

- Logs parseable by Grafana, Datadog, CloudWatch
- Easy to search: `jq '.level == "ERROR"' logs.json`
- Production-ready

---

## Middleware Patterns

### What is Middleware?

**Definition:**

- Function that wraps an HTTP handler
- Executes code BEFORE and/or AFTER the handler
- Can modify request/response or short-circuit execution

**Signature:**

```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // BEFORE handler
        next.ServeHTTP(w, r)
        // AFTER handler
    })
}
```

**Real World:**

- Used by: Every production web application
- Ideal for: Cross-cutting concerns (logging, auth, metrics)

---

### Recovery Middleware (Panic Handling)

**What is it?**

- Catches panics in handlers
- Prevents entire application from crashing
- Returns 500 to client instead of connection reset

**Why critical?**

- ✅ **Availability** - One bug doesn't kill entire app
- ✅ **Debugging** - Logs panic details for investigation
- ✅ **User Experience** - Client gets proper error response

**In our project:**

```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                // Log panic details
                slog.Error("Recovered from panic",
                    "panic", rec,
                    "path", r.URL.Path,
                    "method", r.Method,
                )

                // Return 500 to client
                w.WriteHeader(http.StatusInternalServerError)
            }
        }()

        next.ServeHTTP(w, r)
    })
}
```

**What happens:**

```
Without Recovery:
  User Request → Handler Panic → App Crashes → All users affected ❌

With Recovery:
  User Request → Handler Panic → Recovery catches → 500 response → App continues ✅
```

---

### Logging Middleware (Request/Response Tracking)

**What is it?**

- Logs every HTTP request
- Captures method, path, status code, duration
- Essential for debugging and monitoring

**Why critical?**

- ✅ **Debugging** - See what requests are failing
- ✅ **Monitoring** - Track error rates, slow endpoints
- ✅ **Auditing** - Who accessed what and when

**Problem: Can't read status code from ResponseWriter**

```go
// ❌ This doesn't work!
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        next.ServeHTTP(w, r)

        // How to get status code? w.Status doesn't exist!
        slog.Info("request", "status", ???)
    })
}
```

**Solution: ResponseWriter Wrapper Pattern**

```go
// Wrapper that captures status code
type ResponseRecord struct {
    http.ResponseWriter  // Embed original
    Status  int          // Capture status
    Written bool         // Track if header sent
}

func (r *ResponseRecord) WriteHeader(status int) {
    if r.Written {
        return  // Prevent multiple calls
    }
    r.Status = status
    r.Written = true
    r.ResponseWriter.WriteHeader(status)
}

func (r *ResponseRecord) Write(b []byte) (int, error) {
    if !r.Written {
        r.WriteHeader(http.StatusOK)  // Default 200
    }
    return r.ResponseWriter.Write(b)
}
```

**Usage:**

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap ResponseWriter
        record := &ResponseRecord{
            ResponseWriter: w,
            Status:         http.StatusOK,
        }

        // Pass wrapper to handler
        next.ServeHTTP(record, r)

        // Now we can read status!
        slog.Info("request completed",
            "method", r.Method,
            "path", r.URL.Path,
            "status", record.Status,
            "duration", time.Since(start),
        )
    })
}
```

**How it works:**

```
1. Handler calls w.WriteHeader(201)
2. Go calls ResponseRecord.WriteHeader(201) (our wrapper!)
3. We save: record.Status = 201
4. We call original: ResponseWriter.WriteHeader(201)
5. After handler completes, we can read record.Status
```

---

### Middleware Chaining (Order Matters!)

**In our project:**

```go
handler := middleware.LoggingMiddleware(
    middleware.RecoveryMiddleware(
        mux,
    ),
)
```

**Execution flow:**

```
Request
  ↓
LoggingMiddleware (start timer)
  ↓
RecoveryMiddleware (set defer)
  ↓
Handler (your code)
  ↓
RecoveryMiddleware (check panic)
  ↓
LoggingMiddleware (log status + duration)
  ↓
Response
```

**Why this order?**

- **Logging outermost** - Logs everything (even panics, rate limits)
- **Recovery innermost** - Catches panics from handlers

---

## Security & Compliance

### Rate Limiting (Anti-Scanning Protection)

**Problem: Path Scanning Attacks**

Hackers try hundreds of URLs to find hidden endpoints:

```bash
# Attacker's script:
for path in /admin /api/internal /.env /backup /config; do
    curl https://yourapp.com$path
done
# Result: Hundreds of 404 errors in logs
```

**Solution: Rate Limit 404 Errors**

```go
type IPRateLimiter struct {
    failures map[string]int        // IP → 404 count
    blocked  map[string]time.Time  // IP → block expiry
    mu       sync.RWMutex
}

func RateLimitMiddleware(next http.Handler) http.Handler {
    limiter := NewIPRateLimiter()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr

        // Check if blocked
        if limiter.IsBlocked(ip) {
            w.WriteHeader(http.StatusTooManyRequests)  // 429
            return
        }

        // Execute request
        record := &ResponseRecord{ResponseWriter: w, Status: 200}
        next.ServeHTTP(record, r)

        // Track 404s
        if record.Status == 404 {
            limiter.IncrementFailures(ip)

            if limiter.GetFailures(ip) > 10 {
                limiter.Block(ip, 1*time.Hour)
                slog.Warn("IP blocked", "ip", ip, "reason", "excessive 404s")
            }
        }
    })
}
```

**Effect:**

```
Normal user:
  - Clicks link → 404 (typo) → No problem ✅

Attacker:
  - Request 1-10: 404 responses
  - Request 11+: 429 Too Many Requests (blocked!) ✅
```

**Real World:**

- Used by: Cloudflare, AWS WAF, Nginx rate limiting
- Protects: Against DDoS, scanning, brute force attacks

---

### Audit Logging (Compliance)

**Problem: Regulatory Requirements**

Many industries require logging:

- **GDPR** (Europe) - Who accessed personal data?
- **PCI DSS** (Payments) - Who tried to access card data?
- **HIPAA** (Healthcare) - Who viewed patient records?
- **SOX** (Finance) - Audit trail for all transactions

**Solution: Audit Middleware**

```go
type AuditLog struct {
    Timestamp  time.Time
    IP         string
    UserID     string  // From JWT token
    Method     string
    Path       string
    StatusCode int
    Duration   time.Duration
}

func AuditMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        record := &ResponseRecord{ResponseWriter: w, Status: 200}

        next.ServeHTTP(record, r)

        // Log all requests
        audit := AuditLog{
            Timestamp:  time.Now(),
            IP:         r.RemoteAddr,
            UserID:     getUserID(r),  // From JWT
            Method:     r.Method,
            Path:       r.URL.Path,
            StatusCode: record.Status,
            Duration:   time.Since(start),
        }

        // Structured logging
        slog.Info("audit",
            "timestamp", audit.Timestamp,
            "ip", audit.IP,
            "user_id", audit.UserID,
            "method", audit.Method,
            "path", audit.Path,
            "status", audit.StatusCode,
            "duration", audit.Duration,
        )

        // Alert on unauthorized access
        if record.Status == 401 || record.Status == 403 {
            slog.Warn("unauthorized access attempt",
                "ip", audit.IP,
                "path", audit.Path,
            )
        }
    })
}
```

**Use Cases:**

**1. Security Investigation:**

```bash
# Who tried to access admin panel?
$ grep '"path":"/admin"' audit.log | grep '"status":401'
# Found: 50 attempts from IP 192.168.1.100 → Potential attack!
```

**2. Compliance Audit:**

```bash
# Who deleted event abc-123?
$ grep '"path":"/events/abc-123"' audit.log | grep '"method":"DELETE"'
# Found: user_id="admin-456" at 2026-02-02T10:45:00Z
```

**3. Performance Analysis:**

```bash
# Which endpoints are slow?
$ grep '"duration"' audit.log | jq 'select(.duration > 1000)'
# Found: /events endpoint taking 1.2s
```

**Real World:**

- Required by: Banks, hospitals, e-commerce platforms
- Stored in: Dedicated audit database (long-term retention)
- Monitored by: Security teams, compliance officers

---

## Production Patterns

### Graceful Shutdown

**What is it?**

- Proper cleanup on SIGTERM/SIGINT
- Finish in-flight requests
- Close database connections

**Why Graceful Shutdown?**

- ✅ **Data Integrity** - Don't lose in-flight requests
- ✅ **Clean State** - Close DB connections properly
- ✅ **Zero Downtime** - Kubernetes/Docker expects this

**Real World:**

- Required for: Kubernetes, Docker, Cloud deployments
- Standard pattern in production apps

**In our project:**

```go
func main() {
    srv := &http.Server{Addr: ":8080", Handler: mux}

    // Start server in goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("Server error", "error", err)
        }
    }()

    // Wait for interrupt signal
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)
    <-stop

    slog.Info("Shutting down server...")

    // Graceful shutdown (wait max 20s for requests to finish)
    ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("Shutdown error", "error", err)
    }

    slog.Info("Server stopped")
}
```

**What happens:**

1. SIGTERM received (Ctrl+C, Docker stop, Kubernetes rollout)
2. Server stops accepting new requests
3. Wait for in-flight requests to finish (max 20s)
4. Close database connections
5. Exit cleanly

---

### Error Wrapping (`fmt.Errorf` + `%w`)

**What is it?**

- Go 1.13+ error wrapping
- Preserve error chain
- `errors.Is()` / `errors.As()` checking

**Why Error Wrapping?**

- ✅ **Context** - Add context while preserving original error
- ✅ **Type Checking** - `errors.Is(err, domain.ErrEventNotFound)`
- ✅ **Debugging** - Full error chain in logs

**Real World:**

- Standard in Go 1.13+
- Best practice for error handling

**In our project:**

```go
// ❌ Old way (loses context)
if err != nil {
    return errors.New("failed to create event")
    // Original error lost!
}

// ✅ New way (preserves error chain)
if err != nil {
    return fmt.Errorf("failed to create event: %w", err)
    // Original error preserved!
}

// Check error type
err := repo.GetEvent(ctx, id)
if errors.Is(err, domain.ErrEventNotFound) {
    return 404, "Event not found"
}
```

**Error chain:**

```
failed to run migrations:
  failed to create migration instance:
    failed to open database:
      pq: SSL is not enabled on the server
      ↑ Original error!
```

---

### DTO Pattern (Data Transfer Objects)

**What is it?**

- Separate objects for API requests/responses
- Domain models != API models
- Decoupling

**Why DTOs?**

- ✅ **API Stability** - Change domain without breaking API
- ✅ **Validation** - Validate input before domain
- ✅ **Security** - Don't expose internal fields

**Real World:**

- Standard in REST APIs
- Used by: All production APIs

**In our project:**

```go
// DTO (API layer)
type CreateEventRequest struct {
    Name     string    `json:"name"`
    Price    int       `json:"price"`
    StartAt  time.Time `json:"startAt"`
    EndAt    time.Time `json:"endAt"`
    Capacity int       `json:"capacity"`
}

type EventResponse struct {
    ID             string    `json:"id"`
    Name           string    `json:"name"`
    Price          int       `json:"price"`
    AvailableSpots int       `json:"availableSpots"`
}

// Domain model (internal)
type Event struct {
    id             uuid.UUID  // UUID, not string!
    name           string
    price          int
    availableSpots int
    createdAt      time.Time  // Not exposed in API!
}

// Conversion
func toEventResponse(event *domain.Event) EventResponse {
    return EventResponse{
        ID:             event.ID().String(),  // UUID → string
        Name:           event.Name(),
        Price:          event.Price(),
        AvailableSpots: event.AvailableSpots(),
        // createdAt not included!
    }
}
```

**Benefits:**

- Domain uses UUID, API uses string
- Domain has `createdAt`, API doesn't show it
- Domain changes don't break API contract

---

### Middleware Pattern

**What is it?**

- Wrapper around HTTP handlers
- Cross-cutting concerns (logging, recovery, auth)
- Chainable

**Why Middleware?**

- ✅ **DRY** - Logging logic in one place
- ✅ **Separation** - Cross-cutting concerns separated
- ✅ **Composability** - Stack multiple middlewares

**Real World:**

- Standard in web frameworks
- Used by: Express.js, Django, Spring Boot

**In our project:**

```go
// Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        next.ServeHTTP(w, r)  // Call next handler

        slog.Info("Request",
            "method", r.Method,
            "path", r.URL.Path,
            "duration", time.Since(start),
        )
    })
}

// Recovery middleware (catch panics)
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                slog.Error("Panic recovered", "error", err)
                http.Error(w, "Internal Server Error", 500)
            }
        }()

        next.ServeHTTP(w, r)
    })
}

// Chain middlewares
handler := RecoveryMiddleware(LoggingMiddleware(mux))
```

---

## Key Takeaways

### Performance Choices

| Choice     | Alternative    | Why We Chose It                           |
| ---------- | -------------- | ----------------------------------------- |
| `pgx`      | `database/sql` | 2-3x faster, native Postgres features     |
| `sqlc`     | GORM/Ent       | Zero ORM overhead, full SQL control       |
| `net/http` | Gin/Fiber      | Zero framework overhead, standard library |
| Postgres   | MySQL          | Better concurrency, ACID compliance       |

---

### Architecture Choices

| Pattern            | Benefit                      |
| ------------------ | ---------------------------- |
| Clean Architecture | Testability, maintainability |
| DDD                | Business logic encapsulation |
| Repository Pattern | Database abstraction         |
| DTO Pattern        | API stability                |

---

### Testing Choices

| Tool               | Purpose                         |
| ------------------ | ------------------------------- |
| Testcontainers     | Real database integration tests |
| Table-Driven Tests | Comprehensive test coverage     |
| Race Detector      | Concurrency bug detection       |

---

### Production Readiness

| Feature             | Why                        |
| ------------------- | -------------------------- |
| Graceful Shutdown   | Zero downtime deployments  |
| Structured Logging  | Observability, debugging   |
| Error Wrapping      | Better error context       |
| Embedded Migrations | Self-contained deployments |

---

## Real-World Comparison

**Our stack vs other projects:**

### Startup (MVP)

- ❌ They use: Node.js + MongoDB + Mongoose
- ✅ We use: Go + Postgres + sqlc
- **Why better:** Performance, type safety, scalability

### Mid-size Company

- ❌ They use: Python + Django + ORM
- ✅ We use: Go + Clean Architecture + Repository Pattern
- **Why better:** Concurrency, compile-time safety, faster

### Enterprise

- ✅ They use: Java + Spring Boot + Hibernate
- ✅ We use: Go + DDD + Testcontainers
- **Similar:** Architecture patterns, testing practices
- **Ours better:** Simplicity, deployment (single binary)

---

## Infrastructure & Security (Phase 4 - Planned)

### Redis

**What is it?**

- In-memory data store
- Key-value database
- Extremely fast (microsecond latency)

**Why Redis?**

- ✅ **Speed** - All data in RAM (1M ops/sec)
- ✅ **TTL** - Keys expire automatically (perfect for rate limiting)
- ✅ **Atomic Operations** - INCR, DECR are thread-safe
- ✅ **Distributed** - Works across multiple servers

**Use Cases in GoTicket:**

1. **Rate Limiting**
   ```
   Key: "rate_limit:login:192.168.1.100"
   Value: 5 (attempts)
   TTL: 900 seconds (15 minutes)
   ```

2. **Session Storage** (future)
3. **Caching** (event details, user profiles)

**When to add:** Phase 4 (after auth is complete)

---

### nginx

**What is it?**

- Web server and reverse proxy
- Load balancer
- SSL/TLS terminator

**Why nginx?**

- ✅ **Performance** - Handles 10,000+ concurrent connections
- ✅ **SSL/TLS** - Terminates HTTPS (Let's Encrypt integration)
- ✅ **Load Balancing** - Distributes traffic across multiple Go instances
- ✅ **Static Files** - Serves images, CSS, JS faster than Go
- ✅ **Security** - Rate limiting, DDoS protection

**Use Cases in GoTicket:**

1. **Reverse Proxy**
   ```
   Client → nginx (HTTPS) → Go App (HTTP)
   ```

2. **SSL Termination**
   ```
   nginx handles certificates
   Go app doesn't need to know about SSL
   ```

3. **Load Balancing**
   ```
   nginx → Go Instance 1
         → Go Instance 2
         → Go Instance 3
   ```

**When to add:** Phase 4 (with Docker deployment)

---

### Rate Limiting Strategy

**Problem:** Brute-force attacks on `/auth/login`

**Solution:** Hybrid rate limiting (IP + Email)

**Limits:**
- IP-based: 10 attempts / 15 minutes
- Email-based: 5 attempts / 15 minutes

**Implementation:**
```
Middleware → Check Redis
          → If limit exceeded → 429 Too Many Requests
          → Else → INCR counter → Continue
```

**Why Redis?**
- Atomic INCR (thread-safe)
- Automatic TTL (no cleanup needed)
- Distributed (works with multiple servers)

**When to add:** Phase 4 (requires Redis)

---

## Summary

**This project uses:**

- ✅ **Production-grade tools** (Postgres, Docker, Testcontainers)
- ✅ **Industry-standard patterns** (Clean Architecture, DDD, Repository)
- ✅ **Best practices** (Graceful shutdown, structured logging, error wrapping)
- ✅ **Modern Go** (1.22 routing, slog, embed, generics)
- 📋 **Planned:** Redis (rate limiting), nginx (reverse proxy)

**This is not tutorial code** - this is production-ready architecture used by:

- Uber (Go microservices)
- Netflix (Clean Architecture)
- Spotify (DDD, event-driven)
- Grafana (pgx, sqlc)

**You learned:**

- Backend architecture (Clean + DDD)
- Database best practices (transactions, migrations, pooling)
- Testing strategies (unit, integration, race conditions)
- Production patterns (graceful shutdown, logging, error handling)
- Security (JWT, bcrypt, user enumeration protection, timing attacks)

---

**Congratulations! This is senior-level tech stack!** 🎉
