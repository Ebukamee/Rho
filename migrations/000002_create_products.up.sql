CREATE TABLE products (
    id UUID PRIMARY KEY,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    sku TEXT NOT NULL UNIQUE,
    price BIGINT NOT NULL CHECK (price >= 0),
    currency CHAR(3) NOT NULL,
    image_url TEXT NOT NULL DEFAULT '',
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX products_category_id_idx
    ON products(category_id);

CREATE INDEX products_active_idx
    ON products(active);

CREATE INDEX products_created_at_idx
    ON products(created_at DESC);