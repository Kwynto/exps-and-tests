package main

import (
	"context"
	"log"
	"time"
)

func doWork(ctx context.Context) {
	ctx2, cf := context.WithTimeout(ctx, 30*time.Second)
	defer cf()

	log.Println("starting working...")

	for {
		select {
		case <-ctx2.Done():
			log.Printf("ctx done: %v", ctx2.Err())
			return
		default:
			log.Println("working...")
			time.Sleep(time.Second)
		}
	}
}

func main() {
	ctx, cf := context.WithTimeout(context.Background(), 10*time.Second)
	defer cf()

	doWork(ctx)
}
