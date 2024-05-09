CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    phone_number VARCHAR(16) NOT NULL UNIQUE,
    fullname VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_checkout_id ON checkout_details (checkout_id);
CREATE INDEX IF NOT EXISTS idx_product_id ON checkout_details (product_id);