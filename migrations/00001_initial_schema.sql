-- +goose Up
CREATE TABLE urls (
                      id BIGSERIAL PRIMARY KEY,
                      original_url TEXT NOT NULL,
                      short_code VARCHAR(20) UNIQUE,
                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- idx_urls_short_code allows to instantly find a row by its short_code.
CREATE INDEX idx_urls_short_code ON urls(short_code);

CREATE TABLE clicks (
                        id BIGSERIAL PRIMARY KEY,
                        url_id BIGINT NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
                        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                        user_agent TEXT,
                        ip_address INET
);

-- idx_clicks_url_id_created_at speeds up fetching analytics for a specific URL.
CREATE INDEX idx_clicks_url_id_created_at ON clicks(url_id, created_at DESC);


-- +goose Down
DROP TABLE IF EXISTS clicks;
DROP TABLE IF EXISTS urls;

