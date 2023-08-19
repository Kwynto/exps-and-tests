package main

import (
	"context"
	"fmt"
	"log"
	"sql-exp-test/internal/base"
	"sql-exp-test/internal/storage/mongodb"
	"time"
)

func main() {
	fmt.Println("Start testing MongoDB...")

	mongoHost := "localhost"
	mongoPort := "27017"
	// mongoDB := "entities"

	// urlExample := "mongodb://127.0.0.1:27017"
	path := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

	db, err := mongodb.New(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	base.Run(ctx, db)

	fmt.Println("Test of MongoDB was stoped.")
}
