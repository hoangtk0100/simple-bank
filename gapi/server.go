package gapi

import (
	"fmt"

	db "github.com/hoangtk0100/simple-bank/db/sqlc"
	"github.com/hoangtk0100/simple-bank/pb"
	"github.com/hoangtk0100/simple-bank/token"
	"github.com/hoangtk0100/simple-bank/util"
	"github.com/hoangtk0100/simple-bank/worker"
)

// Server serves gRPC request for banking service
type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
