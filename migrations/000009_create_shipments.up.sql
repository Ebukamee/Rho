CREATE TABLE shipments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE
        REFERENCES orders(id) ON DELETE CASCADE,
    user_id UUID NOT NULL
        REFERENCES users(id) ON DELETE CASCADE,

    carrier VARCHAR(100) NOT NULL,
    service VARCHAR(100) NOT NULL,
    tracking_number VARCHAR(255),

    cost BIGINT NOT NULL
        CHECK (cost >= 0),

    currency CHAR(3) NOT NULL,

    status VARCHAR(30) NOT NULL
        CHECK (
            status IN (
                'pending',
                'processing',
                'shipped',
                'delivered',
                'cancelled'
            )
        ),

    estimated_days INTEGER NOT NULL DEFAULT 0
        CHECK (estimated_days >= 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX shipments_user_id_idx
    ON shipments(user_id);

CREATE INDEX shipments_order_id_idx
    ON shipments(order_id);

CREATE INDEX shipments_tracking_number_idx
    ON shipments(tracking_number);