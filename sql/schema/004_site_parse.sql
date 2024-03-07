-- +goose Up

CREATE TABLE siteParse (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    url_site TEXT NOT NULL,
    type TEXT NOT NULL 
);


-- +goose Down
DROP TABLE siteParse;