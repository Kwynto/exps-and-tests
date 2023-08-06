package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Record struct {
	Name   string `json:"name"`
	Random int    `json:"random"`
}

// var offset int64 = 0

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Need a Kafka topic name.")
		return
	}

	partition := 0
	topic := os.Args[1]
	group := os.Args[2]
	fmt.Println("Kafka topic:", topic)
	fmt.Println("Consumer group:", group)

	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   []string{"localhost:9092"},
			GroupID:   group,
			Topic:     topic,
			Partition: partition,
			MinBytes:  10e3,
			MaxBytes:  10e6,
		})
	// r.SetOffset(offset)

	fmt.Println("Offset start:", r.Offset())

	for {
		// r.SetOffset(offset)
		fmt.Println("Offset 1:", r.Offset())
		m, err := r.ReadMessage(context.Background())
		// m, err := r.FetchMessage(context.Background())
		if err != nil {
			break
		}
		// r.CommitMessages(context.Background(), m)

		// fmt.Println("Offset 2:", r.Offset())

		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		temp := Record{}
		err = json.Unmarshal(m.Value, &temp)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Printf("%T\n", temp)
		// fmt.Println("Offset 3:", r.Offset())
		randTime := random(3, 25)
		time.Sleep(time.Second * time.Duration(randTime))
		// offset++
	}

	r.Close()
}
