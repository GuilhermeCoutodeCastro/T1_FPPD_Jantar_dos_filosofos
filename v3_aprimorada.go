// V3 - Versão Concorrente Aprimorada (T1)
// Baseado no código base filo_conc.go / T1 (Prof. César De Rose)
// Solução por Quebra de Simetria / Ordenação de Recursos (Hierarquia de Garfos)
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
	garfos           [numFilosofos]*Garfo
	temposEspera     [numFilosofos]time.Duration
	maiorEspera      time.Duration
	muMetricas       sync.Mutex
	filosofosComendo int
	maxSimultaneo    int
)

func Filosofo(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	idEsq := id
	idDir := (id + 1) % numFilosofos

	// Quebra da Espera Circular: pega sempre o garfo de menor ID primeiro
	primeiroID, segundoID := idEsq, idDir
	if primeiroID > segundoID {
		primeiroID, segundoID = segundoID, primeiroID
	}

	garfoPrimeiro := garfos[primeiroID]
	garfoSegundo := garfos[segundoID]

	for c := 1; c <= ciclos; c++ {
		// Impressão visual formatada em colunas por filósofo (Pensando)
		for s := 0; s < id*6; s++ {
			fmt.Printf(" ")
		}
		fmt.Printf(" T%d\n", c)
		time.Sleep(tempoPensar)

		inicioEspera := time.Now()

		// Alocação dos garfos em ordem estrita de ID (elimina a espera circular)
		garfoPrimeiro.Lock()
		garfoSegundo.Lock()

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
		garfoSegundo.Unlock()
		garfoPrimeiro.Unlock()
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

	fmt.Printf("\n=== V3 Concorrente Aprimorada Concluída ===\n")
	fmt.Printf("Tempo Total de Execução: %v\n", duracao)
	fmt.Printf("Tempo Médio de Espera por Ciclo: %v\n", mediaEspera)
	fmt.Printf("Maior Tempo de Espera Observado: %v\n", maiorEspera)
	fmt.Printf("Máximo de Filósofos Comendo Simultaneamente: %d\n", maxSimultaneo)
	fmt.Printf("Total de Ciclos Concluídos por Filósofo: %d\n\n", ciclos)
}
