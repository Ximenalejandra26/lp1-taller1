package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Objetivo: Provocar condición de carrera incrementando un contador desde múltiples goroutines,
// luego arreglarla usando Mutex y/o atomic. Ejecuta con el detector de carrera:
//   go run -race ./problema3
// implementa las variantes pedidas.

// Variante insegura (condición de carrera):
func incrementarInseguro(nGoroutines, nIncrementos int) int64 {
	var contador int64 = 0

	var wg sync.WaitGroup
	wg.Add(nGoroutines)

	for i := 0; i < nGoroutines; i++ {
		go func() {
			// asegura wg.Done() se ejecuta al final
			defer wg.Done()

			for j := 0; j < nIncrementos; j++ {
				// incrementar de manera NO atómica (contador = contador + 1)
				contador = contador + 1
			}
		}()
	}

	wg.Wait()
	return contador
}

// Variante con Mutex:
func incrementarConMutex(nGoroutines, nIncrementos int) int64 {
	var contador int64 = 0
	var mu sync.Mutex // Agregamos el Mutex para proteger la variable.
	var wg sync.WaitGroup 
	wg.Add(nGoroutines)

	for i := 0; i < nGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < nIncrementos; j++ {
				//proteger la sección crítica con mu.Lock()/mu.Unlock()
				mu.Lock()
				contador = contador + 1
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return contador
}

// Variante con atomic:
func incrementarConAtomic(nGoroutines, nIncrementos int) int64 {
	var contador int64 = 0
	var wg sync.WaitGroup
	wg.Add(nGoroutines)

	for i := 0; i < nGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < nIncrementos; j++ {
				// usar atomic.AddInt64(&contador, 1)
				atomic.AddInt64(&contador, 1)
			}
		}()
	}

	wg.Wait()
	return contador
}

func main() {
	nGoroutines := 8
	nIncrementos := 100_000

	fmt.Println("=== Variante INSEGURA (espera valor incorrecto, y -race debe reportar):")
	res1 := incrementarInseguro(nGoroutines, nIncrementos)
	fmt.Println("contador:", res1)

	fmt.Println("=== Variante con MUTEX (valor correcto):")
	res2 := incrementarConMutex(nGoroutines, nIncrementos)
	fmt.Println("contador:", res2)

	fmt.Println("=== Variante con ATOMIC (valor correcto):")
	res3 := incrementarConAtomic(nGoroutines, nIncrementos)
	fmt.Println("contador:", res3)

	fmt.Println("Esperado correcto:", int64(nGoroutines*nIncrementos))
}
