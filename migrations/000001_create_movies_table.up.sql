CREATE TABLE IF NOT EXISTS movies
(
    id
    bigserial
    PRIMARY
    KEY,
    created_at
    timestamp(0) with zone NOT NULL DEFAULT NOW(),
    )