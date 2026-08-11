CREATE TABLE inventory (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL UNIQUE
        REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0
        CHECK (quantity >= 0),
    reserved INTEGER NOT NULL DEFAULT 0
        CHECK (reserved >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (reserved <= quantity)
);

CREATE INDEX inventory_product_id_idx
    ON inventory(product_id);