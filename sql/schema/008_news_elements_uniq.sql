-- +goose Up

ALTER TABLE newsElement
ADD CONSTRAINT url_unique UNIQUE (url);

-- +goose Down

ALTER TABLE newsElement
DROP CONSTRAINT IF EXISTS url_unique;