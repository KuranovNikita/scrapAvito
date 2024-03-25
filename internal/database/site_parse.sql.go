// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: site_parse.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSiteParse = `-- name: CreateSiteParse :one
INSERT INTO siteParse (id, created_at, updated_at, name, url_site, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at, updated_at, name, url_site, type, last_fetched_at
`

type CreateSiteParseParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	UrlSite   string
	Type      string
}

func (q *Queries) CreateSiteParse(ctx context.Context, arg CreateSiteParseParams) (Siteparse, error) {
	row := q.db.QueryRowContext(ctx, createSiteParse,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.UrlSite,
		arg.Type,
	)
	var i Siteparse
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.UrlSite,
		&i.Type,
		&i.LastFetchedAt,
	)
	return i, err
}

const getNextSiteParseToFetch = `-- name: GetNextSiteParseToFetch :many
SELECT id, created_at, updated_at, name, url_site, type, last_fetched_at FROM siteParse
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1
`

func (q *Queries) GetNextSiteParseToFetch(ctx context.Context, limit int32) ([]Siteparse, error) {
	rows, err := q.db.QueryContext(ctx, getNextSiteParseToFetch, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Siteparse
	for rows.Next() {
		var i Siteparse
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.UrlSite,
			&i.Type,
			&i.LastFetchedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const markSiteParseAsFetched = `-- name: MarkSiteParseAsFetched :one
UPDATE siteParse
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, name, url_site, type, last_fetched_at
`

func (q *Queries) MarkSiteParseAsFetched(ctx context.Context, id uuid.UUID) (Siteparse, error) {
	row := q.db.QueryRowContext(ctx, markSiteParseAsFetched, id)
	var i Siteparse
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.UrlSite,
		&i.Type,
		&i.LastFetchedAt,
	)
	return i, err
}

const selectSiteParse = `-- name: SelectSiteParse :one
SELECT url_site, type FROM siteParse WHERE id = $1
`

type SelectSiteParseRow struct {
	UrlSite string
	Type    string
}

func (q *Queries) SelectSiteParse(ctx context.Context, id uuid.UUID) (SelectSiteParseRow, error) {
	row := q.db.QueryRowContext(ctx, selectSiteParse, id)
	var i SelectSiteParseRow
	err := row.Scan(&i.UrlSite, &i.Type)
	return i, err
}
