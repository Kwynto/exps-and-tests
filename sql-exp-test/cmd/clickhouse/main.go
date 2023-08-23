package main

import (
	"context"
	"fmt"
	"log"
	"sql-exp-test/internal/base"
	"sql-exp-test/internal/storage/clickhouse"
	"time"
)

func main() {
	fmt.Println("Start testing ClickHouse...")

	chLogin := "default"
	chPassword := ""
	// chPassword = os.Getenv("CLICKHOUSE_PASSWORD")
	chHost := "127.0.0.1"
	chPort := "9000"
	chDB := "entities"

	db, err := clickhouse.New(chHost, chPort, chDB, chLogin, chPassword)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	base.Run(ctx, db)

	fmt.Println("Test of ClickHouse was stoped.")
}
