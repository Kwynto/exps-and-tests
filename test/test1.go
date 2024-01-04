package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.ReadFile("dumbos.txt")
	if err != nil {
		fmt.Println("Unable to open file:", err)
		return
	}
	//fmt.Println(string(file))
	text := string(file)
	//fmt.Println(text)

	for _, v := range text {
		fmt.Println(v)
	}

	// for _, v := range text {
	// 	fmt.Println(string(v))
	// }
}
