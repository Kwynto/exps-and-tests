package main

import (
	"fmt"
	"time"
)

func main() {
	var sl []int = []int{1, 2, 3, 4}

	for i := 0; i < 5; i++ {
		for i := 0; i < 1020; i++ {
			sl = append(sl, i)
			// fmt.Println("Len: ", len(sl), ", Cap: ", cap(sl))
			time.Sleep(time.Millisecond * 5)
		}
		fmt.Println(sl, " Len: ", len(sl), ", Cap: ", cap(sl))

		for i := 0; i < 1000; i++ {
			sl = sl[1:]
			// fmt.Println("Len: ", len(sl), ", Cap: ", cap(sl))
			time.Sleep(time.Millisecond * 5)
		}
		fmt.Println(sl, " Len: ", len(sl), ", Cap: ", cap(sl))
	}

	// sl = nil
	// fmt.Println(sl, " Len: ", len(sl), ", Cap: ", cap(sl))

}
