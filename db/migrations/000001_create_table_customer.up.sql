CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL CHECK (length(name) > 5 AND length(name) <= 50),
    phone_number VARCHAR(30) NOT NULL CHECK (length(phone_number) > 0 AND length(phone_number) <= 30),
    UNIQUE (phone_number)
);

-- Indexes for filtering and sorting
CREATE INDEX IF NOT EXISTS idx_name ON customers (name);
CREATE INDEX IF NOT EXISTS idx_phone_number ON customers (phone_number);

