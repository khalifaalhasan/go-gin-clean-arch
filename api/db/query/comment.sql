-- name: CreateComment :one
INSERT INTO comments (
  username, 
  content, 
  ip_address
) VALUES (
  $1, $2, $3
) RETURNING id, username, content, created_at;

-- name: ListComments :many
SELECT id, username, content, created_at 
FROM comments
WHERE is_hidden = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountCommentsByIP :one
SELECT count(*) 
FROM comments
WHERE ip_address = $1 
AND created_at > (now() - interval '1 minute');