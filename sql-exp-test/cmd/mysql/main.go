package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sql-exp-test/internal/base"
	"sql-exp-test/internal/storage/mysql"
	"time"
)

func main() {
	fmt.Println("Start testing MySQL...")

	mysqlLogin := "root"
	// mysqlPassword := "root"
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlProtocol := "tcp"
	mysqlHost := "localhost"
	mysqlPort := "3306"
	mysqlDB := "entities"

	path := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", mysqlLogin, mysqlPassword, mysqlProtocol, mysqlHost, mysqlPort, mysqlDB)

	db, err := mysql.New(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	base.Run(ctx, db)

	fmt.Println("Test of MySQL was stoped.")
}
