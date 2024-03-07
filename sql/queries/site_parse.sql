-- name: CreateSiteParse :one
INSERT INTO siteParse (id, created_at, updated_at, name, url_site, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

