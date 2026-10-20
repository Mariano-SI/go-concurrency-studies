package main

import "fmt"

func main() {
	// permite o canal receber mais de uam menasgem pra processar ao mesmo tempo

	channel := make(chan string, 2)

	channel <- "hello"
	channel <- "world"

	fmt.Println(<-channel)
	fmt.Println(<-channel)
}
