package main

import (
	"fmt"
	"time"
)

func main() {
	// go fmt.Println("Hello from gorutine.")

	// go fmt.Println("Hello from main.")

	// time.Sleep(time.Millisecond)

	message1 := make(chan string)
	message2 := make(chan string)

	go func() {
		for {
			time.Sleep(time.Millisecond * 500)
			message1 <- "Прошло полсекунды."
		}
	}()

	go func() {
		for {
			time.Sleep(time.Second * 2)
			message2 <- "Прошло 2 секунды."
		}
	}()

	for {
		select {
		case msg := <-message1:
			fmt.Println(msg)
		case msg := <-message2:
			fmt.Println(msg)
		}
	}

}
