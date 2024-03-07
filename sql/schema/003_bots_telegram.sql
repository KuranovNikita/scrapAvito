-- +goose Up

CREATE TABLE botsTelegram (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    bot_token TEXT NOT NULL,
    chat_id TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE    
);


-- +goose Down
DROP TABLE botsTelegram;