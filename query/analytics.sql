-- name: GetClicksByPeriod :many
-- Aggregates click counts for a given URL ID over a specified time period (e.g., 'day', 'month').
SELECT
    date_trunc(sqlc.arg(period)::text, created_at)::date AS key,
    count(*) AS value
FROM clicks
WHERE url_id = sqlc.arg(url_id)
GROUP BY key
ORDER BY key DESC;

-- name: GetClicksByUserAgent :many
-- Aggregates click counts for a given URL ID, grouped by User-Agent.
SELECT
    COALESCE(user_agent, 'Unknown') AS key,
    count(*) as value
FROM clicks
WHERE url_id = sqlc.arg(url_id)
GROUP BY key
ORDER BY value DESC;

-- name: GetClicksByPeriodAndUserAgent :many
-- Aggregates click counts grouped by both a time period AND User-Agent.
SELECT
    date_trunc(sqlc.arg(period)::text, created_at)::date AS time_key,
    COALESCE(user_agent, 'Unknown') AS ua_key,
    count(*) as value
FROM clicks
WHERE url_id = sqlc.arg(url_id)
GROUP BY time_key, ua_key
ORDER BY time_key DESC, value DESC;

