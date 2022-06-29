package gapi

import (
	"context"

	"github.com/croixxant/go-sandbox/simplebank/db"
	"github.com/croixxant/go-sandbox/simplebank/pb"
	"github.com/croixxant/go-sandbox/simplebank/util"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// violations := validateCreateUserRequest(req)
	// if violations != nil {
	// 	return nil, invalidArgumentError(violations)
	// }

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	err = server.store.CreateUser(ctx, arg)
	if err != nil {
		if dbErr, ok := err.(*mysql.MySQLError); ok {
			switch dbErr.Number {
			case 1062:
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

// func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
// 	if err := val.ValidateUsername(req.GetUsername()); err != nil {
// 		violations = append(violations, fieldViolation("username", err))
// 	}

// 	if err := val.ValidatePassword(req.GetPassword()); err != nil {
// 		violations = append(violations, fieldViolation("password", err))
// 	}

// 	if err := val.ValidateFullName(req.GetFullName()); err != nil {
// 		violations = append(violations, fieldViolation("full_name", err))
// 	}

// 	if err := val.ValidateEmail(req.GetEmail()); err != nil {
// 		violations = append(violations, fieldViolation("email", err))
// 	}

// 	return violations
// }
