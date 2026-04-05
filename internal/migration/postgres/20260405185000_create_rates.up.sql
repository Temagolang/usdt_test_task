CREATE TABLE rates (
    id         BIGSERIAL PRIMARY KEY,
    ask        NUMERIC NOT NULL,
    bid        NUMERIC NOT NULL,
    fetched_at TIMESTAMPTZ NOT NULL
);