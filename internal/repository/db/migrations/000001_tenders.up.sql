CREATE TABLE IF NOT EXISTS usdtrub (
    created_at TIMESTAMP PRIMARY KEY DEFAULT CURRENT_TIMESTAMP,
    ask double precision NOT NULL,
    bid double precision NOT NULL
);
