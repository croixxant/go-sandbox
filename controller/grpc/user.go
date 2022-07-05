package grpc

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/croixxant/go-sandbox/controller/grpc/internal"
	"github.com/croixxant/go-sandbox/entity"
	"github.com/croixxant/go-sandbox/usecase"
	"github.com/croixxant/go-sandbox/util"
)

func (c *Controller) CreateUser(ctx context.Context, req *internal.CreateUserRequest) (*internal.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	args := usecase.CreateUserParams{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
	u, err := usecase.NewUserUsecase(c.repo, c.tokenmaker, c.config).CreateUser(ctx, args)
	if err != nil {
		if dbErr, ok := err.(*mysql.MySQLError); ok {
			switch dbErr.Number {
			case 1062:
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	rsp := &internal.CreateUserResponse{
		User: convertUser(u),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *internal.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := util.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := util.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}

func (c *Controller) LoginUser(ctx context.Context, req *internal.LoginUserRequest) (*internal.LoginUserResponse, error) {
	violations := validateLoginUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	mtdt := c.extractMetadata(ctx)
	args := usecase.LoginUserParams{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		UserAgent: mtdt.UserAgent,
		ClientIP:  mtdt.ClientIP,
	}
	r, err := usecase.NewUserUsecase(c.repo, c.tokenmaker, c.config).LoginUser(ctx, args)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	rsp := &internal.LoginUserResponse{
		User:                  convertUser(r.User),
		SessionId:             r.SessionID.String(),
		AccessToken:           r.AccessToken,
		RefreshToken:          r.RefreshToken,
		AccessTokenExpiresAt:  timestamppb.New(r.AccessTokenExpiresAt),
		RefreshTokenExpiresAt: timestamppb.New(r.RefreshTokenExpiresAt),
	}
	return rsp, nil
}

func validateLoginUserRequest(req *internal.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := util.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := util.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}

func convertUser(user *entity.User) *internal.User {
	return &internal.User{
		Username:  user.Username,
		FullName:  user.FullName,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
