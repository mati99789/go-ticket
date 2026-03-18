# Load Tests — k6

## Purpose

Verify that the booking system **does not double-book** under real HTTP load (full stack: middleware, auth, rate limiting, database).

This complements the integration tests (`TestEventRepository_ReserveSpots_Concurrent`) which only test the repository layer in isolation.

---

## Files

| File | Status | Description |
| --------------------- | ------------ | ----------------------------------------------------- |
| `seed.sql` | Ready | Test data: organizer user + event (capacity=10000) |
| `booking_scenario.js` | Complete | k6 scenario: setup() login → default() concurrent bookings |

---

## Test Data (seed.sql)

The seed creates:

- **User**: `organizer@loadtest.com` / `LoadTest123!` / role: `organizer`
- **Event**: `a0000000-0000-0000-0000-000000000001` / capacity: `10000`

bcrypt hash (cost=12): `$2a$12$YwRyt8wAQrp2TX9rNDLSye9TZ2pgILU7cMVbi4ecYAqA6EsIPRQge`

---

## How to Run

### 1. Install k6

```bash
sudo snap install k6
# or via apt:
sudo apt-get install gnupg
curl -s https://dl.k6.io/key.gpg | sudo apt-key add -
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6
k6 version
```

### 2. Start the application

```bash
docker compose up -d
```

### 3. Apply seed data

```bash
docker compose exec -T db psql -U postgres -d goticket -f /dev/stdin < tests/load/seed.sql
```

Verify seed was applied:

```bash
docker compose exec db psql -U postgres -d goticket -c "SELECT email, role FROM users;"
docker compose exec db psql -U postgres -d goticket -c "SELECT name, available_spots FROM events;"
```

### 4. Run the load test

```bash
k6 run tests/load/booking_scenario.js
```

### 5. Verify — check for double-bookings after the test

```bash
docker compose exec db psql -U postgres -d goticket \
  -c "SELECT COUNT(*) FROM bookings WHERE event_id = 'a0000000-0000-0000-0000-000000000001';"

docker compose exec db psql -U postgres -d goticket \
  -c "SELECT available_spots FROM events WHERE id = 'a0000000-0000-0000-0000-000000000001';"
```

**Expected result**: `COUNT(*) + available_spots = 10000`

---

## Scenario: booking_scenario.js

### Thresholds

| Metric | Threshold | Reason |
| ------------------------- | --------- | ------------------------ |
| `http_req_duration` p(95) | < 200ms | SLA for the booking API |
| `http_req_failed` | < 1% | Acceptable error rate |

### Stages (ramp-up)

```
10s  → 10 VUs   (warm-up)
30s  → 50 VUs   (normal traffic)
20s  → 100 VUs  (peak load)
10s  → 0 VUs    (ramp-down)
```

### Scenario Structure

```javascript
import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = { stages: [...], thresholds: {...} }

export function setup() {
    // Runs ONCE before the test
    // POST /auth/login with seed credentials
    // Returns { token, eventId } shared across all VUs
}

export default function(data) {
    // Runs for every VU on every iteration
    // POST /events/{data.eventId}/bookings
    // Authorization: Bearer {data.token}
    // Expect: 201 (success) or 409 (event full)
}
```

---

## Proven Result

**Verified**: 1031 bookings created + 8969 remaining spots = 10000 total — zero double-bookings under 100 concurrent users.

This proves the `UPDATE events SET available_spots = available_spots - 1 WHERE available_spots > 0` atomic pattern works correctly under real concurrent HTTP load.

---

## CI/CD Integration (Future)

```yaml
# .github/workflows/load-test.yml
- name: Seed test data
  run: docker compose exec -T db psql -U postgres -d goticket -f /dev/stdin < tests/load/seed.sql

- name: Run k6 load test
  uses: grafana/k6-action@v0.3.0
  with:
    filename: tests/load/booking_scenario.js
```
