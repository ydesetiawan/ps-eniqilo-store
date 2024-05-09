CREATE TABLE checkout_details (
    id SERIAL PRIMARY KEY,
    checkout_id SERIAL NOT NULL,
    product_id SERIAL NOT NULL,
    product_price NUMERIC NOT NULL CHECK (product_price >= 1),
    total_price NUMERIC NOT NULL CHECK (total_price >= 1),
    quantity INTEGER NOT NULL CHECK (quantity >= 1),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_checkout_id ON checkout_details (checkout_id);
CREATE INDEX IF NOT EXISTS idx_product_id ON checkout_details (product_id);