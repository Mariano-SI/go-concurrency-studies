package main

import (
	"fmt"
	"sync"
	"time"
)

func task(taskName string, wg *sync.WaitGroup) {

	defer wg.Done()
	for i := 0; i < 5; i++ {
		fmt.Printf("%d: Task %s is running\n", i, taskName)
		time.Sleep(1 * time.Second)
	}
}

func main() {

	fmt.Println("Exemplo com go routines: execução concorrente")
	start := time.Now()
	waitGroup := sync.WaitGroup{}

	waitGroup.Add(2)
	go task("A - With goroutines", &waitGroup)
	go task("B - With goroutines", &waitGroup)

	waitGroup.Wait()
	fmt.Printf("Tempo total com goroutines: %s\n", time.Since(start))

}
