package main

import (
	"fmt"
	"time"
)

func task(taskName string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("%d: Task %s is running\n", i, taskName)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	fmt.Println("Exemplo sem go routines: execução sequencial")
	start := time.Now()
	task("A - Without goroutines")
	task("B - Without goroutines")
	fmt.Printf("Tempo total sem goroutines: %s\n", time.Since(start))

	fmt.Println("Exemplo com go routines: execução concorrente")
	start = time.Now()
	go task("A - With goroutines")
	go task("B - With goroutines")

	// ATENÇÃO: usar time.Sleep para segurar a main NÃO é a forma correta de
	// esperar goroutines terminarem. Aqui funciona só porque sabemos que cada
	// task leva ~5s. Se uma task demorar mais que o sleep, a main encerra antes
	// e a goroutine é interrompida; se demorar menos, ficamos esperando à toa.
	// A forma correta é usar sync.WaitGroup (ou channels).
	time.Sleep(6 * time.Second)

	fmt.Printf("Tempo total com goroutines: %s\n", time.Since(start))

}
