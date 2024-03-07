-- name: CreateBotTelegram :one
INSERT INTO botsTelegram (id, created_at, updated_at, name, bot_token, chat_id, user_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;
 