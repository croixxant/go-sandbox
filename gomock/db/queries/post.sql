-- name: GetUserPosts :many
SELECT * FROM posts
WHERE user_id = ?;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = ?;

-- name: AddPostLikesCount :exec
UPDATE posts
SET likes_count = likes_count + sqlc.arg(likes)
WHERE id = sqlc.arg(id);

-- name: CreatePost :execresult
INSERT INTO posts (
  user_id, body, likes_count
) VALUES (
  ?, ?, 0
);

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = ?;
