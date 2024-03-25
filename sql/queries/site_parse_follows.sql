-- name: CreateSiteParseFollows :one
INSERT INTO siteParseFollows (id, created_at, updated_at, user_id, site_parse_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetSiteParseFollows :many
SELECT user_id, site_parse_id FROM siteParseFollows;

-- name: GetBotDataBySiteParseID :many
SELECT b.bot_token, b.chat_id
FROM siteParseFollows spf
JOIN botsTelegram b ON spf.user_id = b.user_id
WHERE spf.site_parse_id = $1;