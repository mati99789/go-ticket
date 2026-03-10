-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (id, event_name, event_data, destination, aggregate_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetPendingOutboxEvents :many
SELECT * FROM outbox_events
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkOutBoxEventAsProcessed :exec
UPDATE outbox_events
SET status = 'processed', updated_at = NOW()
WHERE id = $1;