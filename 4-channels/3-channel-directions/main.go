package main

import "fmt"

// main goroutine
func main() {
	channel := make(chan int)

	go func() {
		channel <- 1000 // canal recebe dado
	}()

	msg := <-channel // canal esvazia

	fmt.Println(msg)

	go publish(channel)
	receiver(channel)
}

func publish(ch chan<- int) { //publisher é um canal que recebe mensagens
	for i := 0; i < 10; i++ {
		fmt.Println("[publisher] - enviando:", i)
		ch <- i
	}
	close(ch)
}

func receiver(ch <-chan int) { // receiver é um canal que le mensagens
	for x := range ch {
		fmt.Println("[reader] - processando:", x)
	}
}
