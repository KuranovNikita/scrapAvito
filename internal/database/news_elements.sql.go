// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: news_elements.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const clearNewsElement = `-- name: ClearNewsElement :exec
DELETE FROM newsElement
`

func (q *Queries) ClearNewsElement(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, clearNewsElement)
	return err
}

const createNewsElement = `-- name: CreateNewsElement :one
INSERT INTO newsElement (id, created_at, updated_at, site_parse_id, title, news_date, url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (url) DO NOTHING
RETURNING id, created_at, updated_at, site_parse_id, title, news_date, url
`

type CreateNewsElementParams struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	SiteParseID uuid.UUID
	Title       string
	NewsDate    string
	Url         string
}

func (q *Queries) CreateNewsElement(ctx context.Context, arg CreateNewsElementParams) (Newselement, error) {
	row := q.db.QueryRowContext(ctx, createNewsElement,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.SiteParseID,
		arg.Title,
		arg.NewsDate,
		arg.Url,
	)
	var i Newselement
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.SiteParseID,
		&i.Title,
		&i.NewsDate,
		&i.Url,
	)
	return i, err
}