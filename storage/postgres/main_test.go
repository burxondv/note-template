package postgres_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/burxondv/note-template/config"
	"github.com/burxondv/note-template/storage"
	"github.com/jmoiron/sqlx"
)

var (
	strg storage.StorageI
)

func TestMain(m *testing.M) {
	cfg := config.Load("./../..")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to open connection: %v", err)
	}

	strg = storage.NewStoragePg(db)
	os.Exit(m.Run())
}
