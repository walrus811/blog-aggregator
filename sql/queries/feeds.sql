-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: DeleteFeedById :exec
DELETE FROM feeds WHERE id = $1;

-- name: GetNextFeedsToFetch :many
select * from feeds 
where last_fetched_at is null or last_fetched_at < NOW() - interval '1 days'  
order by last_fetched_at nulls first
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET last_fetched_at = NOW(), updated_at = NOW() WHERE id = $1
RETURNING *;