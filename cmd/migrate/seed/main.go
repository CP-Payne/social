package main

import (
	"log"

	"github.com/CP-Payne/social/internal/db"
	"github.com/CP-Payne/social/internal/env"
	"github.com/CP-Payne/social/internal/store"
)

func main() {

	addr := env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/socialnetwork?sslmode=disable")

	conn, err := db.New(addr, 15, 15, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
