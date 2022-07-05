package repo

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/croixxant/go-sandbox/entity"
)

// Repository usecaseからの依存は実装ではなくこちらに向ける
// 集約で分割するのが望ましいだろうとはなる
type Repository interface {
	AddWalletBalance(ctx context.Context, arg AddWalletBalanceParams) error
	CreateWallet(ctx context.Context, arg CreateWalletParams) (int64, error)
	DeleteWallet(ctx context.Context, id int64) error
	GetEntry(ctx context.Context, id int64) (*entity.Entry, error)
	GetTransfer(ctx context.Context, id int64) (*entity.Transfer, error)
	GetWallet(ctx context.Context, id int64) (*entity.Wallet, error)
	ListEntries(ctx context.Context, arg ListEntriesParams) ([]*entity.Entry, error)
	ListTransfers(ctx context.Context, arg ListTransfersParams) ([]*entity.Transfer, error)
	ListWallets(ctx context.Context, arg ListwalletsParams) ([]*entity.Wallet, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (*CreateTransferResult, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (int64, error)
	GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error)
	GetUser(ctx context.Context, username string) (*entity.User, error)
}

type CreateEntryParams struct {
	WalletID int64 `json:"wallet_id"`
	Amount   int64 `json:"amount"`
}

type ListEntriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListTransfersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type AddWalletBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

type CreateWalletParams struct {
	UserID   int64  `json:"user_id"`
	Balance  int64  `json:"balance"`
	Currency string `json:"currency"`
}

type ListwalletsParams struct {
	UserID int64 `json:"user_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type CreateTransferParams struct {
	FromWalletID int64 `json:"from_wallet_id"`
	ToWalletID   int64 `json:"to_wallet_id"`
	Amount       int64 `json:"amount"`
}

type CreateTransferResult struct {
	TransferID  int64 `json:"transfer_id"`
	FromEntryID int64 `json:"from_entry"`
	ToEntryID   int64 `json:"to_entry"`
}

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type CreateUserParams struct {
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}
