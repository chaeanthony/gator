-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id) 
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsByUserName :many
SELECT * FROM feeds WHERE user_id=(SELECT id FROM users WHERE users.name=$1);

-- name: GetFeeds :many
SELECT *, (SELECT name FROM users WHERE id=feeds.user_id) FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds 
SET last_fetched_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP 
WHERE id = $1 RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds ORDER BY last_fetched_at NULLS FIRST LIMIT 1;



