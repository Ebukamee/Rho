CREATE TABLE payments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(id),
    user_id UUID NOT NULL REFERENCES users(id),
    provider VARCHAR(50) NOT NULL,
    provider_ref VARCHAR(255) NOT NULL UNIQUE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency CHAR(3) NOT NULL,
    status VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CHECK (
        status IN (
            'pending',
            'succeeded',
            'failed',
            'refunded'
        )
    )
);

CREATE INDEX payments_order_idx
    ON payments(order_id);

CREATE INDEX payments_user_idx
    ON payments(user_id);

CREATE INDEX payments_provider_idx
    ON payments(provider);