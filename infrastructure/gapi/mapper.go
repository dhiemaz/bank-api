package gapi

import (
	"github.com/dhiemaz/bank-api/grpc/pb"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromDBUserToPbUserResponse(user db.User) *pb.UserResponse {
	return &pb.UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
	}
}
