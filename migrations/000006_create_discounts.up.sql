CREATE TABLE discounts (
    id UUID PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(20) NOT NULL,
    value BIGINT NOT NULL CHECK (value > 0),
    minimum_order BIGINT NOT NULL DEFAULT 0
        CHECK (minimum_order >= 0),
    usage_limit INTEGER
        CHECK (usage_limit IS NULL OR usage_limit > 0),
    usage_count INTEGER NOT NULL DEFAULT 0
        CHECK (usage_count >= 0),
    starts_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (type IN ('percentage', 'fixed')),
    CHECK (
        type != 'percentage'
        OR value <= 100
    ),
    CHECK (
        expires_at IS NULL
        OR starts_at IS NULL
        OR expires_at > starts_at
    )
);

CREATE INDEX discounts_code_idx
    ON discounts(code);

CREATE INDEX discounts_active_idx
    ON discounts(active);