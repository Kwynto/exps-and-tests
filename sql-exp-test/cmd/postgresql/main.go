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

	pgsqlLogin := "root"
	// pgsqlPassword := "root"
	pgsqlPassword := os.Getenv("POSTGRESQL_PASSWORD")
	pgsqlProtocol := "tcp"
	pgsqlHost := "localhost"
	pgsqlPort := "3306"
	pgsqlDB := "entities"

	path := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", pgsqlLogin, pgsqlPassword, pgsqlProtocol, pgsqlHost, pgsqlPort, pgsqlDB)

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
