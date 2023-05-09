package main

import (
	"database/sql"
	"golang-backend-master/api"
	db "golang-backend-master/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:1234@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = ":8888"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannnot start server:", err)
	}
}
