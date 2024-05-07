CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL CHECK (length(name) > 0 AND length(name) <= 30),
    sku VARCHAR(30) NOT NULL CHECK (length(sku) > 0 AND length(sku) <= 30),
    category VARCHAR(20) NOT NULL CHECK (category IN ('Clothing', 'Accessories', 'Footwear', 'Beverages')),
    image_url TEXT NOT NULL CHECK (image_url ~* '^https?://'),
    notes VARCHAR(200) NOT NULL CHECK (length(notes) > 0 AND length(notes) <= 200),
    price NUMERIC NOT NULL CHECK (price >= 1),
    stock INTEGER NOT NULL CHECK (stock >= 0 AND stock <= 100000),
    location VARCHAR(200) NOT NULL CHECK (length(location) > 0 AND length(location) <= 200),
    is_available BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexes for filtering and sorting
CREATE INDEX IF NOT EXISTS idx_name ON products (name);
CREATE INDEX IF NOT EXISTS idx_sku ON products (sku);
CREATE INDEX IF NOT EXISTS idx_category ON products (category);
CREATE INDEX IF NOT EXISTS idx_price ON products (price);
CREATE INDEX IF NOT EXISTS idx_created_at ON products (created_at);