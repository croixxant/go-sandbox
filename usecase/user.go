package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/croixxant/go-sandbox/config"
	"github.com/croixxant/go-sandbox/entity"
	"github.com/croixxant/go-sandbox/usecase/repo"
	"github.com/croixxant/go-sandbox/util"
	"github.com/croixxant/go-sandbox/util/token"
)

type UserUsecase struct {
	repo       repo.Repository
	tokenmaker token.Maker
	config     config.Config
}

func NewUserUsecase(r repo.Repository, tm token.Maker, cfg config.Config) *UserUsecase {
	return &UserUsecase{
		repo:       r,
		tokenmaker: tm,
		config:     cfg,
	}
}

type CreateUserParams struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserUsecase) CreateUser(ctx context.Context, args CreateUserParams) (*entity.User, error) {
	hashedPassword, err := util.HashPassword(args.Password)
	if err != nil {
		return nil, err
	}

	rArg := repo.CreateUserParams{
		Username:       args.Username,
		FullName:       args.FullName,
		Email:          args.Email,
		HashedPassword: hashedPassword,
	}
	_, err = u.repo.CreateUser(ctx, rArg)
	if err != nil {
		return nil, err
	}

	return u.repo.GetUser(ctx, args.Username)
}

type LoginUserParams struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

type LoginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  *entity.User `json:"user"`
}

func (u *UserUsecase) LoginUser(ctx context.Context, args LoginUserParams) (*LoginUserResponse, error) {
	usr, err := u.repo.GetUser(ctx, args.Username)
	if err != nil {
		return nil, err
	}

	err = util.CheckPassword(args.Password, usr.HashedPassword)
	if err != nil {
		return nil, err
	}

	aToken, aPayload, err := u.tokenmaker.CreateToken(
		usr.ID,
		u.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	rToken, rPayload, err := u.tokenmaker.CreateToken(
		usr.ID,
		u.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	err = u.repo.CreateSession(ctx, repo.CreateSessionParams{
		ID:           rPayload.ID,
		UserID:       usr.ID,
		RefreshToken: rToken,
		UserAgent:    args.UserAgent,
		ClientIp:     args.ClientIP,
		IsBlocked:    false,
		ExpiresAt:    rPayload.ExpiredAt,
	})
	if err != nil {
		return nil, err
	}

	resp := &LoginUserResponse{
		SessionID:             rPayload.ID,
		AccessToken:           aToken,
		AccessTokenExpiresAt:  aPayload.ExpiredAt,
		RefreshToken:          rToken,
		RefreshTokenExpiresAt: rPayload.ExpiredAt,
		User:                  usr,
	}
	return resp, nil
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (u *UserUsecase) RenewAccessToken(ctx context.Context, args RenewAccessTokenRequest) (*RenewAccessTokenResponse, error) {
	rPayload, err := u.tokenmaker.VerifyToken(args.RefreshToken)
	if err != nil {
		return nil, err
	}

	sess, err := u.repo.GetSession(ctx, rPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	if sess.IsBlocked {
		err := fmt.Errorf("blocked session")
		return nil, err
	}

	if sess.UserID != rPayload.UserID {
		err := fmt.Errorf("incorrect session user")
		return nil, err
	}

	if sess.RefreshToken != args.RefreshToken {
		err := fmt.Errorf("mismatched session token")
		return nil, err
	}

	if time.Now().After(sess.ExpiresAt) {
		err := fmt.Errorf("expired session")
		return nil, err
	}

	aToken, aPayload, err := u.tokenmaker.CreateToken(
		rPayload.UserID,
		u.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	resp := &RenewAccessTokenResponse{
		AccessToken:          aToken,
		AccessTokenExpiresAt: aPayload.ExpiredAt,
	}
	return resp, nil
}
