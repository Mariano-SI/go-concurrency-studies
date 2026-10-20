package main

import "fmt"

// routine main
func main() {
	channel := make(chan string) // canal vazio

	go func() {
		channel <- "Olá, Mundo!" // outra rotoutine preencheu o canal // canal cheio
	}()

	msg := <-channel // esvaziou canal

	fmt.Println(msg)
}
