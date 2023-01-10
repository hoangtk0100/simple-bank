package gapi

import (
	"time"

	db "github.com/hoangtk0100/simple-bank/db/sqlc"
	"github.com/hoangtk0100/simple-bank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: convertTimestamp(user.PasswordChangedAt),
		CreatedAt:         convertTimestamp(user.CreatedAt),
	}
}

func convertTimestamp(input time.Time) *timestamppb.Timestamp {
	return timestamppb.New(input)
}
