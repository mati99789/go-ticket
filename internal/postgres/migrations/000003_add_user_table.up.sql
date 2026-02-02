CREATE TYPE user_role AS ENUM ('user', 'admin', 'organizer');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);


CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role ON users (role);