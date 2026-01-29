-- name: CreateEvent :one
INSERT INTO events (id, name, price, start_at, end_at, created_at, updated_at, capacity, available_spots)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateEvent :one
UPDATE events
SET name = $2, price = $3, start_at = $4, end_at = $5, updated_at = $6, capacity = $7
WHERE id = $1
RETURNING *;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;

-- name: GetEvent :one
SELECT * FROM events
WHERE id = $1;

-- name: ListEvents :many
SELECT * FROM events
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: ReserveSpots :one
UPDATE events
SET available_spots = available_spots - $2
WHERE id = $1 AND available_spots >= $2
RETURNING *;
