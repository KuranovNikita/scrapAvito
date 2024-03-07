-- +goose Up

ALTER TABLE siteParse
ADD CONSTRAINT unique_name_url_site UNIQUE (name, url_site);

-- +goose Down

ALTER TABLE siteParse
DROP CONSTRAINT IF EXISTS unique_name_url_site;