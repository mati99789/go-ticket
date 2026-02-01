DROP TABLE bookings;
ALTER TABLE events DROP COLUMN capacity;
ALTER TABLE events DROP COLUMN available_spots;

DROP INDEX idx_bookings_event_id;
DROP INDEX idx_bookings_user_email;
