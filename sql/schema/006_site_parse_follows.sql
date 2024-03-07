-- +goose Up

CREATE TABLE siteParseFollows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    site_parse_id UUID NOT NULL REFERENCES siteParse(id) ON DELETE CASCADE,
    UNIQUE(user_id, site_parse_id) 
);


-- +goose Down
DROP TABLE siteParseFollows;