-- name: CreateNewsElement :one
INSERT INTO newsElement (id, created_at, updated_at, site_parse_id, title, news_date, url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url) DO NOTHING
RETURNING *;


-- name: ClearNewsElement :exec
DELETE FROM newsElement;
