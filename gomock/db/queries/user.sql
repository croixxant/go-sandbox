-- name: GetUser :one
SELECT * FROM users
WHERE id = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :execresult
INSERT INTO users (
  email, hashed_password
) VALUES (
  ?, ?
);

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: AddUserLikesCount :exec
UPDATE users
SET likes_count = likes_count + sqlc.arg(likes)
WHERE id = sqlc.arg(id);
