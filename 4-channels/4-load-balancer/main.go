package main

import (
	"fmt"
	"time"
)

func worker(workerId int, data <-chan int) {
	for x := range data {
		fmt.Printf("Worker %d recebeu %d\n", workerId, x)
		time.Sleep(time.Second)
	}
}

func main() {

	data := make(chan int)

	go worker(1, data)
	go worker(2, data)
	go worker(3, data)

	for i := 0; i < 100; i++ {
		data <- i
	}

}
