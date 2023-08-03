package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wg sync.WaitGroup

func reader(ch chan string) {
	for {
		v := <-ch
		fmt.Printf("v: %v\n", v)
	}
}

func go1() {
	for {
		fmt.Println("go1")
		time.Sleep(300 * time.Millisecond)
	}
}

func go2() {
	for {
		fmt.Println("go2")
		time.Sleep(300 * time.Millisecond)
	}
}

func go3() {
	for {
		fmt.Println("go3")
		time.Sleep(300 * time.Millisecond)
	}
}

func gomain() {
	fmt.Println("Main gorutine was starting.")
	go go1()
	go go2()
	go go3()
	fmt.Println("Main gorutine was stoping.")
	wg.Done()
}

func getGoMaxProc() int {
	return runtime.GOMAXPROCS(0)
}

func main() {

	a1 := "A-1"
	a2 := "A-2"

	for {
		break
	}

	// a2 = for {
	// 	return "ret A-2"
	// }

	a2, a1 = a1, a2

	fmt.Println(a1)
	fmt.Println(a2)

	// var intChannel chan int

	// intChannel := make(chan string)
	// go reader(intChannel)
	// intChannel <- "a"
	// fmt.Println("main")
	// intChannel <- "b"
	// intChannel <- "c"
	// intChannel <- "d"
	// intChannel <- "e"
	// time.Sleep(10 * time.Millisecond)

	// wg.Add(1)
	// go gomain()
	// wg.Wait()
	// time.Sleep(5 * time.Second)
	// fmt.Println("Programm was stoped.")

	// fmt.Printf("GOMAXPROCS: %d\n", getGoMaxProc())

	_, err := fmt.Printf("GOMAXPROCS: %d\n", getGoMaxProc())

	if err != nil {
		fmt.Println("Блин!")
	}

	// cmd := exec.Command("ifconfig", "-a")
	// cmd := exec.Command("netstat", "-nr")
	// b, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// fmt.Println("Result: ", string(b))

	// var arr = [...]int{1, 2, 3, 4} // Массив
	var arr = []int{1, 2, 3, 4} // Слайс

	pArr := &arr
	pElArr := &arr[1]

	fmt.Println("Массив: ", arr)
	fmt.Println("Указатель: ", pArr)
	fmt.Println("Указатель на элемент:", pElArr)
	fmt.Println("Значение указателя:", *pElArr)

}
