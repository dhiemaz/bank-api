package db

import (
	"database/sql"
	"github.com/dhiemaz/bank-api/config"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
	
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestMain(m *testing.M) {
	var err error
	// Load config
	config := config.GetConfig()
	if err != nil {
		log.Fatal("cannot load configuration for testing", err)
	}

	// Connect to db
	testDB, err = sql.Open(config.Database.Driver, config.Database.URL)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	// Set connection & run tests
	testQueries = New(testDB)
	os.Exit(m.Run())
}
