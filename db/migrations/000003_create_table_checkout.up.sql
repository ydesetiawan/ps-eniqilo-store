CREATE TABLE checkouts (
    id SERIAL PRIMARY KEY,
    customer_id SERIAL NOT NULL,
    total_price NUMERIC NOT NULL CHECK (total_price >= 1),
    paid NUMERIC NOT NULL CHECK (paid >= 1),
    change NUMERIC NOT NULL CHECK (change >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX IF NOT EXISTS idx_cid ON checkouts (customer_id);
CREATE INDEX IF NOT EXISTS idx_checkout_created_at ON checkouts (created_at);