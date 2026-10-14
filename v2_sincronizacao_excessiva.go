// v2_sincronizacao_excessiva.go (V2 – Versão com Sincronização Excessiva)
// V2 - Versão sem Deadlock com Sincronização Excessiva (T1)
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

var (
	arbitro          sync.Mutex // Árbitro global: causa a sincronização excessiva
	garfos           [numFilosofos]*Garfo
	temposEspera     [numFilosofos]time.Duration
	maiorEspera      time.Duration
	muMetricas       sync.Mutex
	filosofosComendo int
	maxSimultaneo    int
)

func Filosofo(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	garfoEsq := garfos[id]
	garfoDir := garfos[(id+1)%numFilosofos]

	for c := 1; c <= ciclos; c++ {
		// Impressão visual formatada em colunas por filósofo (Pensando)
		for s := 0; s < id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" T%d\n", c)
		time.Sleep(tempoPensar)

		inicioEspera := time.Now()

		// Sincronização Excessiva: Apenas um filósofo por vez pode negociar a posse dos garfos
		arbitro.Lock()
		garfoEsq.Lock()
		garfoDir.Lock()
		arbitro.Unlock()

		espera := time.Since(inicioEspera)

		// Coleta de métricas de espera e concorrência
		muMetricas.Lock()
		temposEspera[id] += espera
		if espera > maiorEspera {
			maiorEspera = espera
		}
		filosofosComendo++
		if filosofosComendo > maxSimultaneo {
			maxSimultaneo = filosofosComendo
		}
		muMetricas.Unlock()

		// Impressão visual formatada em colunas por filósofo (Comendo)
		for s := 0; s < id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" E%d\n", c)
		time.Sleep(tempoComer)

		muMetricas.Lock()
		filosofosComendo--
		muMetricas.Unlock()

		// Liberação dos garfos
		garfoDir.Unlock()
		garfoEsq.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < numFilosofos; i++ {
		garfos[i] = &Garfo{}
	}

	fmt.Println("\n [P0] [P1] [P2] [P3] [P4]")
	inicio := time.Now()

	for i := 0; i < numFilosofos; i++ {
		wg.Add(1)
		go Filosofo(i, &wg)
	}

	wg.Wait()
	duracao := time.Since(inicio)

	// Cálculo do tempo médio de espera
	var totalEspera time.Duration
	for _, t := range temposEspera {
		totalEspera += t
	}
	mediaEspera := totalEspera / time.Duration(numFilosofos*ciclos)

	fmt.Printf("\n=== V2 Sincronização Excessiva Concluída ===\n")
	fmt.Printf("Tempo Total de Execução: %v\n", duracao)
	fmt.Printf("Tempo Médio de Espera por Ciclo: %v\n", mediaEspera)
	fmt.Printf("Maior Tempo de Espera Observado: %v\n", maiorEspera)
	fmt.Printf("Máximo de Filósofos Comendo Simultaneamente: %d\n", maxSimultaneo)
	fmt.Printf("Total de Ciclos Concluídos por Filósofo: %d\n\n", ciclos)
}
