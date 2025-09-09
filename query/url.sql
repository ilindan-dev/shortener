-- name: CreateURL :one
-- Inserts a new URL record with the original URL.
INSERT INTO urls (original_url)
VALUES ($1)
RETURNING *;

-- name: UpdateURLShortCode :exec
-- Updates a URL record with its generated short code.
UPDATE urls
SET short_code = $2
WHERE id = $1;

-- name: GetURLByShortCode :one
-- Retrieves a URL record by its unique short code.
SELECT *
FROM urls
WHERE short_code = $1;

-- name: CreateClick :exec
-- Inserts a new click record for analytics.
INSERT INTO clicks (url_id, user_agent, ip_address)
VALUES ($1, $2, $3);

-- name: GetClicksByURLID :many
-- Retrieves all click records for a given URL, ordered by the most recent.
SELECT *
FROM clicks
WHERE url_id = $1
ORDER BY created_at DESC;