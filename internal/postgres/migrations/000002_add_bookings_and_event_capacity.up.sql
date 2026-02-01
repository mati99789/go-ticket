ALTER TABLE events ADD COLUMN capacity INT NOT NULL DEFAULT 0;
ALTER TABLE events ADD COLUMN available_spots INT NOT NULL DEFAULT 0;

CREATE TABLE bookings (
    id UUID PRIMARY KEY NOT NULL,
    event_id UUID NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

CREATE INDEX idx_bookings_event_id ON bookings(event_id);
CREATE INDEX idx_bookings_user_email ON bookings(user_email);
