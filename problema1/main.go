package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Objetivo: Lanzar varias goroutines que imprimen mensajes y esperar a que todas terminen.

// Completa los pasos marcados con TODO para entender goroutines y WaitGroup.
func worker(id int, veces int, wg *sync.WaitGroup) {
	// asegurar que al finalizar la función se haga wg.Done()
	
	defer wg.Done()

	for i := 1; i <= veces; i++ {
		fmt.Printf("[worker %d] hola %d\n", id, i)
		// dormir un poco para simular trabajo (p. ej. 100–300 ms)
		time.Sleep(time.Duration(100+rand.Intn(200)) * time.Millisecond)
		
	}
}

func main() {
	var wg sync.WaitGroup

	// cambiar estos parámetros y observar el intercalado de salidas
	numGoroutines := 5 
	veces := 3

	// lanzar varias goroutines, sumar al WG y esperar con wg.Wait()
	for id := 1; id <= numGoroutines; id++ {
		wg.Add(1)
		go worker(id, veces, &wg)
	}

	// Esperar a que todas las goroutines terminen
	wg.Wait()

	fmt.Println("Listo: todas las goroutines terminaron.")
}
