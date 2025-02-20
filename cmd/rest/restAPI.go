package rest

import (
	"fmt"
	"github.com/dhiemaz/bank-api/config"
	"github.com/dhiemaz/bank-api/infrastructure"
	"github.com/dhiemaz/bank-api/infrastructure/db/sqlc"
	"log"
)

//	@title			Bank API
//	@version		1.0
//	@description	Bank API is a SAAP that allows users to create accounts and transfer money between them.
//
//	@contact.email	prawira.dimas.yudha@gmail.com
//	@contact.name	Dimas Yudha Prawira

// @securityDefinitions.apikey	bearerAuth
// @in							header
// @name						Authorization
// @description				Bearer <token>
func Run() {
	// Load config from environment variables
	config := config.GetConfig()

	// Initialize the database
	conn := db.InitDatabase(config)
	store := db.NewStore(conn)
	query := db.New(conn)

	ginServer, err := infrastructure.NewServer(config, store, query)
	if err != nil {
		log.Fatalf("cannot create HTTP server, err: %s", err)
	}

	port := "8000"
	log.Printf("HTTP server is listening on %s", port)
	if err := ginServer.Start(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
		log.Fatalf("cannot start server address: %s, err: %s", port, err)
	}
}
