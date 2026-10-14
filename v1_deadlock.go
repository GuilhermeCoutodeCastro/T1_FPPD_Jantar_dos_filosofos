// V1 - Versão Concorrente Sujeita a Deadlock (T1)
// Baseado no código base filo_conc.go / T1 (Prof. César De Rose)
package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	numFilosofos = 5
	ciclos       = 20
	tempoPensar  = 2 * time.Millisecond
	tempoComer   = 3 * time.Millisecond
)

type Garfo struct {
	sync.Mutex
}

func Filosofo(id int, garfos []*Garfo, wg *sync.WaitGroup) {
	defer wg.Done()

	garfoEsq := garfos[id]
	garfoDir := garfos[(id+1)%numFilosofos]

	for c := 1; c <= ciclos; c++ {
		// Impressão visual formatada em colunas por filósofo
		for s := 0; s < id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" T%d\n", c)

		// Pensando
		time.Sleep(tempoPensar)

		// 1. Tenta pegar o garfo da esquerda
		garfoEsq.Lock()

		// Pausa intencional para induzir o entrelaçamento e o deadlock em todos
		time.Sleep(1 * time.Millisecond)

		// 2. Tenta pegar o garfo da direita (causa bloqueio mútuo)
		garfoDir.Lock()

		for s := 0; s < id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" E%d\n", c)

		// Comendo
		time.Sleep(tempoComer)

		garfoDir.Unlock()
		garfoEsq.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	garfos := make([]*Garfo, numFilosofos)

	for i := 0; i < numFilosofos; i++ {
		garfos[i] = &Garfo{}
	}

	fmt.Println("\n [P0] [P1] [P2] [P3] [P4]")

	// Dispara cada filósofo como uma goroutine concorrente
	for i := 0; i < numFilosofos; i++ {
		wg.Add(1)
		go Filosofo(i, garfos, &wg)
	}

	wg.Wait()
}
