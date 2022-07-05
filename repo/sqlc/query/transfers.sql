-- name: CreateTransfer :execlastid
INSERT INTO transfers (
  from_wallet_id,
  to_wallet_id,
  amount
) VALUES (
  ?, ?, ?
);

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = ? LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
ORDER BY id
LIMIT ?
OFFSET ?;
