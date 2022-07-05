package usecase

import (
	"github.com/croixxant/go-sandbox/usecase/repo"
)

type Usecase struct {
	repo repo.Repository
}

func NewUsecase(r repo.Repository) *Usecase {
	return &Usecase{
		repo: r,
	}
}
