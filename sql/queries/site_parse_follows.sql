-- name: CreateSiteParseFollows :one
INSERT INTO siteParseFollows (id, created_at, updated_at, user_id, site_parse_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;
