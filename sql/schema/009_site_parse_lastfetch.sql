-- +goose Up

ALTER TABLE siteParse ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE siteParse DROP COLUMN last_fetched_at;