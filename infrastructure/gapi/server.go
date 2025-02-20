package gapi

import (
	"fmt"
	"github.com/dhiemaz/bank-api/config"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/utils/token"
	"log"
	"net"

	"github.com/dhiemaz/bank-api/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	config *config.Config
	db     db.Store
	token  token.Maker
	pb.UnimplementedBankServiceServer
}

func NewServer(config *config.Config, store db.Store) (*GRPCServer, error) {
	maker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker for grpcServer, %w", err)
	}

	grpcServer := &GRPCServer{config: config, token: maker, db: store}
	return grpcServer, nil
}

func (server *GRPCServer) Start(address string) error {
	grpcServer := grpc.NewServer()
	pb.RegisterBankServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening on %s", address)
	return nil
}
