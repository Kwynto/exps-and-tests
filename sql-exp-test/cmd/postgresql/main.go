package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sql-exp-test/internal/base"
	"sql-exp-test/internal/storage/postgresql"
	"time"
)

func main() {
	fmt.Println("Start testing PostgreSQL...")

	pgsqlLogin := "postgres"
	// pgsqlPassword := "postgres"
	pgsqlPassword := os.Getenv("POSTGRESQL_PASSWORD")
	// pgsqlProtocol := "tcp"
	pgsqlHost := "localhost"
	pgsqlPort := "5432"
	pgsqlDB := "entities"
	sslmode := "disable"

	// urlExample := "postgres://username:password@localhost:5432/database_name?sslmode=disable"
	path := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", pgsqlLogin, pgsqlPassword, pgsqlHost, pgsqlPort, pgsqlDB, sslmode)

	db, err := postgresql.New(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	base.Run(ctx, db)

	fmt.Println("Test of PostgreSQL was stoped.")
}
