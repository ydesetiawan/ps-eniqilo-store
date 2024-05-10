CREATE TABLE users (
    id bigserial PRIMARY KEY,
    name VARCHAR(50),
    phone_number VARCHAR(30) NOT NULL CHECK (length(phone_number) > 0 AND length(phone_number) <= 30),
    password VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_name ON customers (name);
CREATE INDEX IF NOT EXISTS idx_user_phone_number ON customers (phone_number);
