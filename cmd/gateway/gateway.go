package gateway

import (
	"context"
	"github.com/dhiemaz/bank-api/config"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"github.com/dhiemaz/bank-api/infrastructure/gapi"
	"log"
	"net"
	"net/http"

	"github.com/dhiemaz/bank-api/grpc/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func RunGateway() {
	// Load config from environment variables
	config := config.GetConfig()

	// Initialize the database
	conn := db.InitDatabase(config)
	store := db.NewStore(conn)

	grpcServer, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("cannot create gRPC server, err: %s", err)
	}

	jsonOpts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOpts)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := pb.RegisterBankServiceHandlerServer(ctx, grpcMux, grpcServer); err != nil {
		log.Fatalf("cannot register gRPC server, err %s", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	setupSwagger(mux, config)

	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalf("cannot listen on port 8000, err: %s", err)
	}

	log.Printf("Gateway server is listening on 8000")
	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("cannot start HTTP server, err: %s", err)
	}
}

func setupSwagger(mux *http.ServeMux, config *config.Config) {
	swaggerFileHandler := http.FileServer(http.Dir("./docs/swagger"))
	mux.Handle("/docs/", http.StripPrefix("/docs/", swaggerFileHandler))
}
