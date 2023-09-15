package main

import (
	"database/sql"
	"log"

	"github.com/anilbolat/simple-bank/util"

	"github.com/anilbolat/simple-bank/api"
	db "github.com/anilbolat/simple-bank/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("error while loading the config file.")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
