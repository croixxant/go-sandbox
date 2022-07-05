package usecase

import (
	"context"
	"errors"

	"github.com/croixxant/go-sandbox/entity"
	"github.com/croixxant/go-sandbox/usecase/repo"
)

func (u *Usecase) CreateWallet(ctx context.Context, args repo.CreateWalletParams) (*entity.Wallet, error) {
	id, err := u.repo.CreateWallet(ctx, args)
	if err != nil {
		return nil, err
	}

	return u.repo.GetWallet(ctx, id)
}

type GetWalletParams struct {
	RequestUserID int64
	ID            int64
}

func (u *Usecase) GetWallet(ctx context.Context, args GetWalletParams) (*entity.Wallet, error) {
	w, err := u.repo.GetWallet(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if w.UserID != args.RequestUserID {
		err := errors.New("wallet doesn't belong to the authenticated user")
		return nil, err
	}
	return w, nil
}

func (u *Usecase) ListWallets(ctx context.Context, args repo.ListwalletsParams) ([]*entity.Wallet, error) {
	return u.repo.ListWallets(ctx, args)
}
