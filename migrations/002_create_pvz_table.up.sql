CREATE TABLE IF NOT EXISTS pvz
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    city       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);