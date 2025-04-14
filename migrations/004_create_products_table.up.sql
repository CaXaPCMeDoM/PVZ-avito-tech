CREATE TABLE IF NOT EXISTS products
(
    id           UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    reception_id UUID      NOT NULL REFERENCES receptions (id) ON DELETE CASCADE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    type         VARCHAR(255) NOT NULL
);