// V0 - Versão Sequencial (Baseline para T1)
// Baseado no codigo base filo_seq.go (Prof. Cesar De Rose)
package main

import (
	"fmt"
	"time"
)

func Philosopher(id, iteration int, timeToThink, timeToEat time.Duration) {
	// Imprime recuo visual para identificar a coluna de cada filósofo
	for s := 0; s < id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" T%d\n", iteration)

	// Pensando
	time.Sleep(timeToThink)

	for s := 0; s < id*6; s++ {
		fmt.Printf(" ")
	}
	fmt.Printf(" E%d\n", iteration)

	// Comendo
	time.Sleep(timeToEat)
}

func main() {
	philosophers := 5
	rounds := 20 // 20 ciclos conforme o roteiro do experimento

	// Tempos simulados para pensar e comer
	timeToThink := 2 * time.Millisecond
	timeToEat := 3 * time.Millisecond

	fmt.Println("\n [P0] [P1] [P2] [P3] [P4]")
	start := time.Now()

	// Execução sequencial dos filósofos por ciclo
	for i := 1; i <= rounds; i++ {
		for j := 0; j < philosophers; j++ {
			Philosopher(j, i, timeToThink, timeToEat)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("\n=== V0 Sequencial Concluído ===\n")
	fmt.Printf("Tempo Total de Execução: %s\n", elapsed)
	fmt.Printf("Total de Ciclos por Filósofo: %d\n", rounds)
	fmt.Printf("Máximo de Filósofos Comendo Simultaneamente: 1\n\n")
}
