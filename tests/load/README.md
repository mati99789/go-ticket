# Load Tests — k6

## Cel

Weryfikacja że system rezerwacji **nie double-bookuje** pod prawdziwym obciążeniem HTTP (cały stack: middleware, auth, rate limiting, DB).

Uzupełnienie do integration testów (`TestEventRepository_ReserveSpots_Concurrent`) które testują tylko warstwę repozytorium.

## Pliki

| Plik                  | Status       | Opis                                                  |
| --------------------- | ------------ | ----------------------------------------------------- |
| `seed.sql`            | ✅ Gotowy    | Dane testowe: organizer user + event (capacity=10000) |
| `booking_scenario.js` | 🚧 W trakcie | Scenariusz k6 — wymaga dokończenia (patrz niżej)      |

## Dane testowe (seed.sql)

Seed tworzy:

- **User**: `organizer@loadtest.com` / `LoadTest123!` / rola: `organizer`
- **Event**: `a0000000-0000-0000-0000-000000000001` / capacity: `10000`

Hasło bcrypt (cost=12): `$2a$12$YwRyt8wAQrp2TX9rNDLSye9TZ2pgILU7cMVbi4ecYAqA6EsIPRQge`

## Jak uruchomić (docelowo)

### 1. Zainstaluj k6

```bash
sudo snap install k6
# lub przez apt:
sudo apt-get install gnupg
curl -s https://dl.k6.io/key.gpg | sudo apt-key add -
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6
k6 version
```

### 2. Uruchom aplikację

```bash
docker-compose up -d
```

### 3. Zaaplikuj seed data

```bash
docker-compose exec -T db psql -U postgres -d goticket -f /dev/stdin < tests/load/seed.sql
```

Weryfikacja:

```bash
docker-compose exec db psql -U postgres -d goticket -c "SELECT email, role FROM users;"
docker-compose exec db psql -U postgres -d goticket -c "SELECT name, available_spots FROM events;"
```

### 4. Uruchom test

```bash
k6 run tests/load/booking_scenario.js
```

### 5. Weryfikacja po teście — czy nie ma double-bookingu

```bash
docker-compose exec db psql -U postgres -d goticket \
  -c "SELECT COUNT(*) FROM bookings WHERE event_id = 'a0000000-0000-0000-0000-000000000001';"

docker-compose exec db psql -U postgres -d goticket \
  -c "SELECT available_spots FROM events WHERE id = 'a0000000-0000-0000-0000-000000000001';"
```

Oczekiwany wynik: `COUNT(*) + available_spots = 10000`

## Scenariusz (booking_scenario.js)

### Thresholds

| Metryka                   | Threshold | Powód                    |
| ------------------------- | --------- | ------------------------ |
| `http_req_duration` p(95) | < 200ms   | SLA dla API bookingowego |
| `http_req_failed`         | < 1%      | Dopuszczalny error rate  |

### Stages (ramp-up)

```
10s  → 10 VU   (warm-up)
30s  → 50 VU   (normalny ruch)
20s  → 100 VU  (peak)
10s  → 0 VU    (ramp-down)
```

### ⚠️ Co jeszcze wymaga implementacji w booking_scenario.js

Obecny stan pliku ma `options` i `default()` — **brakuje**:

#### 1. Funkcja `setup()` — uruchamia się RAZ przed testem

```
Pseudokod:
export function setup() {
    // 1. POST /auth/login z danymi z seed.sql
    //    body: { email: "organizer@loadtest.com", password: "LoadTest123!" }
    //    Oczekiwany response: { token: "eyJ..." }

    // 2. Wyciągnij token z JSON response
    //    const token = JSON.parse(res.body).token

    // 3. Zwróć dane dla VU
    //    return { token: token, eventId: 'a0000000-0000-0000-0000-000000000001' }
}
```

#### 2. Zaktualizuj `default(data)` — przyjmuje dane z `setup()`

```
Pseudokod:
export default function(data) {
    // payload: { userEmail: "organizer@loadtest.com" }
    // URL: /events/{data.eventId}/bookings
    // Header: Authorization: Bearer {data.token}
    // Oczekiwany status: 201 (lub 409 gdy event full)
}
```

#### 3. Usuń hardkodowany event_id z payload body

Endpoint `CreateBooking` oczekuje w body tylko `userEmail` — event_id idzie w URL path.

### Docelowa struktura pliku

```javascript
import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = { stages: [...], thresholds: {...} }

export function setup() {
    // login → return { token, eventId }
}

export default function(data) {
    // POST /events/{data.eventId}/bookings
    // Authorization: Bearer {data.token}
    // body: { userEmail: "..." }
}
```

## Integracja z CI/CD (przyszłość)

```yaml
# .github/workflows/load-test.yml
- name: Seed test data
  run: docker-compose exec -T db psql ... < tests/load/seed.sql

- name: Run k6
  uses: grafana/k6-action@v0.3.0
  with:
    filename: tests/load/booking_scenario.js
```
