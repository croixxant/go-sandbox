package sqlc

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang/mock/mockgen/model" // https://github.com/golang/mock#debugging-errors
	"github.com/google/uuid"

	"github.com/croixxant/go-sandbox/entity"
	"github.com/croixxant/go-sandbox/repo/sqlc/internal"
	"github.com/croixxant/go-sandbox/usecase/repo"
)

// Repository usecase.Repositoryの実装。
// モデルを受け取ってなんちゃらとかそういう感じにはしていない
type Repository struct {
	q  internal.Querier
	db *sql.DB
}

// callされていなくともinterfaceを実装していなければコンパイルエラーが起こらないようにしている
var _ repo.Repository = (*Repository)(nil)

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		q:  internal.New(db),
		db: db,
	}
}

func (s *Repository) tx(ctx context.Context, fn func(*internal.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := internal.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func (s *Repository) CreateTransfer(ctx context.Context, arg repo.CreateTransferParams) (*repo.CreateTransferResult, error) {
	var result repo.CreateTransferResult

	transfer := func(q *internal.Queries) error {
		transferID, err := q.CreateTransfer(ctx, internal.CreateTransferParams(arg))
		if err != nil {
			return err
		}
		result.TransferID = transferID

		fromEntryID, err := q.CreateEntry(ctx, internal.CreateEntryParams{
			WalletID: arg.FromWalletID,
			Amount:   -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntryID = fromEntryID

		toEntryID, err := q.CreateEntry(ctx, internal.CreateEntryParams{
			WalletID: arg.ToWalletID,
			Amount:   arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntryID = toEntryID

		if arg.FromWalletID < arg.ToWalletID {
			err = addMoney(ctx, q, arg.FromWalletID, -arg.Amount, arg.ToWalletID, arg.Amount)
		} else {
			err = addMoney(ctx, q, arg.ToWalletID, arg.Amount, arg.FromWalletID, -arg.Amount)
		}
		return err
	}
	err := s.tx(ctx, transfer)

	return &result, err
}

func (s *Repository) AddWalletBalance(ctx context.Context, arg repo.AddWalletBalanceParams) error {
	return s.q.AddWalletBalance(ctx, internal.AddWalletBalanceParams(arg))
}
func (s *Repository) CreateWallet(ctx context.Context, arg repo.CreateWalletParams) (int64, error) {
	return s.q.CreateWallet(ctx, internal.CreateWalletParams(arg))
}
func (s *Repository) DeleteWallet(ctx context.Context, id int64) error {
	return s.q.DeleteWallet(ctx, id)
}
func (s *Repository) GetEntry(ctx context.Context, id int64) (*entity.Entry, error) {
	e, err := s.q.GetEntry(ctx, id)
	if err != nil {
		return nil, err
	}
	return (*entity.Entry)(e), nil
}
func (s *Repository) GetTransfer(ctx context.Context, id int64) (*entity.Transfer, error) {
	t, err := s.q.GetTransfer(ctx, id)
	if err != nil {
		return nil, err
	}
	return (*entity.Transfer)(t), nil
}
func (s *Repository) GetWallet(ctx context.Context, id int64) (*entity.Wallet, error) {
	w, err := s.q.GetWallet(ctx, id)
	if err != nil {
		return nil, err
	}
	return (*entity.Wallet)(w), nil
}
func (s *Repository) ListEntries(ctx context.Context, arg repo.ListEntriesParams) ([]*entity.Entry, error) {
	l, err := s.q.ListEntries(ctx, internal.ListEntriesParams(arg))
	if err != nil {
		return nil, err
	}
	es := make([]*entity.Entry, 0, len(l))
	for _, v := range l {
		es = append(es, (*entity.Entry)(v))
	}
	return es, nil
}
func (s *Repository) ListTransfers(ctx context.Context, arg repo.ListTransfersParams) ([]*entity.Transfer, error) {
	l, err := s.q.ListTransfers(ctx, internal.ListTransfersParams(arg))
	if err != nil {
		return nil, err
	}
	ts := make([]*entity.Transfer, 0, len(l))
	for _, v := range l {
		ts = append(ts, (*entity.Transfer)(v))
	}
	return ts, nil
}
func (s *Repository) ListWallets(ctx context.Context, arg repo.ListwalletsParams) ([]*entity.Wallet, error) {
	l, err := s.q.ListWallets(ctx, internal.ListWalletsParams(arg))
	if err != nil {
		return nil, err
	}
	ws := make([]*entity.Wallet, 0, len(l))
	for _, v := range l {
		ws = append(ws, (*entity.Wallet)(v))
	}
	return ws, nil
}
func (s *Repository) CreateSession(ctx context.Context, arg repo.CreateSessionParams) error {
	return s.q.CreateSession(ctx, internal.CreateSessionParams(arg))
}
func (s *Repository) CreateUser(ctx context.Context, arg repo.CreateUserParams) (int64, error) {
	return s.q.CreateUser(ctx, internal.CreateUserParams(arg))
}
func (s *Repository) GetSession(ctx context.Context, id uuid.UUID) (*entity.Session, error) {
	sess, err := s.q.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}
	return (*entity.Session)(sess), nil
}
func (s *Repository) GetUser(ctx context.Context, username string) (*entity.User, error) {
	u, err := s.q.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return (*entity.User)(u), nil
}

func addMoney(ctx context.Context, q *internal.Queries, walletID1 int64, amount1 int64, walletID2 int64, amount2 int64) (err error) {
	err = q.AddWalletBalance(ctx, internal.AddWalletBalanceParams{
		ID:     walletID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	err = q.AddWalletBalance(ctx, internal.AddWalletBalanceParams{
		ID:     walletID2,
		Amount: amount2,
	})
	return
}
