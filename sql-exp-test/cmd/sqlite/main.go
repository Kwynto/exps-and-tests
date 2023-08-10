package main

import (
	"context"
	"fmt"
	"log"
	"sql-exp-test/internal/base"
	"sql-exp-test/internal/storage/sqlite"
	"time"
)

func main() {
	fmt.Println("Start testing SQLite3...")

	db, err := sqlite.New("./data/sqlite.db")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	base.Run(ctx, db)

	fmt.Println("Test of SQLite3 was stoped.")
}
