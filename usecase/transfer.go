package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/croixxant/go-sandbox/entity"
	"github.com/croixxant/go-sandbox/usecase/repo"
)

type CreateTransferParams struct {
	RequestUserID int64
	FromWalletID  int64
	ToWalletID    int64
	Amount        int64
	Currency      string
}

func (u *Usecase) CreateTransfer(ctx context.Context, arg CreateTransferParams) (*repo.CreateTransferResult, error) {
	fromWallet, err := u.validWallet(ctx, arg.FromWalletID, arg.Currency)
	if err != nil {
		return nil, err
	}
	if fromWallet.UserID != arg.RequestUserID {
		err := errors.New("from wallet doesn't belong to the authenticated user")
		return nil, err
	}
	_, err = u.validWallet(ctx, arg.ToWalletID, arg.Currency)
	if err != nil {
		return nil, err
	}

	rArg := repo.CreateTransferParams{
		FromWalletID: arg.FromWalletID,
		ToWalletID:   arg.ToWalletID,
		Amount:       arg.Amount,
	}
	return u.repo.CreateTransfer(ctx, rArg)
}

func (u *Usecase) validWallet(ctx context.Context, walletD int64, currency string) (*entity.Wallet, error) {
	w, err := u.repo.GetWallet(ctx, walletD)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if w.Currency != currency {
		err := fmt.Errorf("w [%d] currency mismatch: %s vs %s", w.ID, w.Currency, currency)
		return nil, err
	}

	return w, nil
}
