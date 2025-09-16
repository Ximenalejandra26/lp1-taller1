package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Objetivo: Simular "futuros" en Go usando canales. Una función lanza trabajo asíncrono
// y retorna un canal de solo lectura con el resultado futuro.
// completa las funciones y experimenta con varios futuros a la vez.

func asyncCuadrado(x int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		// simular trabajo
		tiempo := time.Duration(rand.Intn(500)+50) * time.Millisecond
		time.Sleep(tiempo)
		ch <- x * x
	}()
	return ch
}

// fanIn combina múltiples canales de entrada en un solo canal de salida.
// Esta función debe estar en el nivel de paquete, fuera de main.
func fanIn(in ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)
	wg.Add(len(in))
	for _, ch := range in {
		go func(c <-chan int) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// crea varios futuros y recolecta sus resultados: f1, f2, f3
	fmt.Println("Lanzando futuros...")
	f1 := asyncCuadrado(2)
	f2 := asyncCuadrado(3)
	f3 := asyncCuadrado(4)

	//Opción 1: esperar cada futuro secuencialmente
	fmt.Println("--- Esperando resultados secuencialmente ---")
	res1 := <-f1
	fmt.Printf("Resultado 1: %d\n", res1)
	res2 := <-f2
	fmt.Printf("Resultado 2: %d\n", res2)
	res3 := <-f3
	fmt.Printf("Resultado 3: %d\n", res3)

	fmt.Println("Todos los futuros terminaron de forma secuencial.")

	// Opción 2: fan-in (combinar múltiples canales)
	fmt.Println("--- Esperando resultados con fan-in ---")
	f4 := asyncCuadrado(5)
	f5 := asyncCuadrado(6)
	f6 := asyncCuadrado(7)

	fanInChan := fanIn(f4, f5, f6)
	
	for i := 0; i < 3; i++ {
		res := <-fanInChan
		fmt.Printf("Resultado Fan-In: %d\n", res)
	}

}