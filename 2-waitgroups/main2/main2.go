package main

import (
	"fmt"
	"sync"
	"time"
)

// A task não precisa mais receber o *sync.WaitGroup nem chamar wg.Done():
// quem cuida da contagem é o wg.Go.
func task(taskName string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d: Task %s is running\n", i, taskName)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	fmt.Println("Exemplo com wg.Go (Go 1.25+)")
	start := time.Now()
	var waitGroup sync.WaitGroup

	// wg.Go(f) faz internamente o equivalente a:
	//   waitGroup.Add(1)
	//   go func() {
	//       defer waitGroup.Done()
	//       f()
	//   }()
	// Por isso não é preciso informar a quantidade com Add(n).
	waitGroup.Go(func() { task("A - With wg.Go") })
	waitGroup.Go(func() { task("B - With wg.Go") })

	waitGroup.Wait()
	fmt.Printf("Tempo total com wg.Go: %s\n", time.Since(start))
}
