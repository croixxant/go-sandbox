-- name: CreateEntry :execlastid
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  ?, ?
);

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = ? LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
LIMIT ?
OFFSET ?;
