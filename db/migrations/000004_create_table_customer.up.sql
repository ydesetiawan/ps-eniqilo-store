CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    fullname VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_phone ON customers (phone_number);
CREATE INDEX IF NOT EXISTS idx_name ON customers (fullname);