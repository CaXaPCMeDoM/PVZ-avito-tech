CREATE TABLE IF NOT EXISTS receptions
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    pvz_id     UUID      NOT NULL REFERENCES pvz (id) ON DELETE CASCADE,
    status     VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);