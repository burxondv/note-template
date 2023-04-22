package main

import (
	"fmt"
	"log"

	"github.com/burxondv/note-template/api"
	"github.com/burxondv/note-template/config"
	"github.com/burxondv/note-template/storage"
	"github.com/go-redis/redis/v9"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	apiServer := api.New(&api.RouterOptions{
		Cfg:      &cfg,
		Storage:  strg,
		InMemory: inMemory,
	})

	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	log.Print("server stopped")
}
