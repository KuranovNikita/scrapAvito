-- name: CreateSiteParse :one
INSERT INTO siteParse (id, created_at, updated_at, name, url_site, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: SelectSiteParse :one
SELECT url_site, type FROM siteParse WHERE id = $1;

-- name: GetNextSiteParseToFetch :many
SELECT * FROM siteParse
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkSiteParseAsFetched :one
UPDATE siteParse
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;