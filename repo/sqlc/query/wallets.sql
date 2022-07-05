-- name: CreateWallet :execlastid
INSERT INTO wallets (
  user_id,
  balance,
  currency
) VALUES (
  ?, ?, ?
);

-- name: GetWallet :one
SELECT * FROM wallets
WHERE id = ? LIMIT 1;

-- -- name: GetWalletForUpdate :one
-- SELECT * FROM wallets
-- WHERE id = ? LIMIT 1
-- FOR UPDATE;

-- name: ListWallets :many
SELECT * FROM wallets
WHERE user_id = ?
ORDER BY id
LIMIT ?
OFFSET ?;

-- -- name: UpdateWallet :exec
-- UPDATE wallets
-- SET balance = ?
-- WHERE id = ?;

-- name: AddWalletBalance :exec
UPDATE wallets
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id);

-- name: DeleteWallet :exec
DELETE FROM wallets
WHERE id = ?;