-- name: CreateBooking :one
INSERT INTO bookings (id, event_id, user_email, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateBooking :one
UPDATE bookings
SET event_id = $2, user_email = $3, status = $4, updated_at = $5
WHERE id = $1
RETURNING *;

-- name: DeleteBooking :exec
DELETE FROM bookings
WHERE id = $1;

-- name: GetBookingByID :one
SELECT * FROM bookings
WHERE id = $1;

-- name: ListBookings :many
SELECT * FROM bookings
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: ConfirmBooking :one
UPDATE bookings
SET status = 'confirmed', updated_at = NOW()
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: CancelBooking :one
UPDATE bookings
SET status = 'cancelled', updated_at = NOW()
WHERE id = $1 AND status = 'pending'
RETURNING *;
