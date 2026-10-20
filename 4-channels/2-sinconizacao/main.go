package main

import "fmt"

func main() {
	channel := make(chan int)

	go publish(channel)

	for x := range channel {
		fmt.Println("[reader] - processando:", x)
	}
}

func publish(ch chan int) {
	for i := 0; i < 10; i++ {
		fmt.Println("[publisher] - enviando:", i)
		ch <- i
	}
	close(ch)
}

