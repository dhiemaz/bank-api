package gapi

import (
	"github.com/dhiemaz/bank-api/config"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/infrastructure/gapi"
	"log"
)

func RunGRPCAPIServer() {
	// Load config from environment variables
	config := config.GetConfig()

	// Initialize the database
	conn := db.InitDatabase(config)
	store := db.NewStore(conn)
	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	log.Printf("GRPC server is listening on %s", "8000")
	if err := grpcServer.Start("0.0.0.0:8000"); err != nil {
		log.Fatalf("cannot start gRPC server address: %s, err: %s", "8000", err)
	}
}
