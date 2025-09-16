package main

import (
	"fmt"
	"sync"
	"time"
)

// Objetivo: Provocar un deadlock con dos mutex y dos goroutines que adquieren
// recursos en orden distinto. Luego evitarlo imponiendo un orden global.
// NOTA: La versión con deadlock se quedará bloqueada: ejecútala, observa y luego cambia a la versión segura.
// completa/activa la sección que quieras probar.

func deadlock() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("G1: Lock mu1") 
		// adquirir mu1
		mu1.Lock()

		time.Sleep(100 * time.Millisecond) // fuerza entrelazado
		fmt.Println("G1: Lock mu2") 
		// adquirir mu2
		mu2.Lock()

		fmt.Println("G1: listo")
	}()

	go func() {
		defer wg.Done()
		fmt.Println("G2: Lock mu2") 
		// adquirir mu2
		mu2.Lock()

		time.Sleep(100 * time.Millisecond)
		fmt.Println("G2: Lock mu1") 
		// adquirir mu1
		mu1.Lock()

		fmt.Println("G2: listo")
	}()

	// ADVERTENCIA: esto no retornará por el deadlock
	wg.Wait()
}

func seguroOrdenado() {
	var mu1, mu2 sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	// Regla: siempre adquirir mu1 luego mu2 (mismo orden en todos)
	lockEnOrden := func(a, b *sync.Mutex) func() func() {
		// retorna: lock():unlock()
		return func() func() {
			// adquirir a luego b
			a.Lock()
			b.Lock()

			// retorna función para liberar en orden inverso
			return func() {
				// liberar b luego a
				b.Unlock()
				a.Unlock()
			}
		}
	}

	go func() {
		defer wg.Done()
		unlock := lockEnOrden(&mu1, &mu2)()
		defer unlock()
		fmt.Println("G1: trabajo con mu1->mu2")
		time.Sleep(100 * time.Millisecond)
	}()

	go func() {
		defer wg.Done()
		unlock := lockEnOrden(&mu1, &mu2)()
		defer unlock()
		fmt.Println("G2: trabajo con mu1->mu2")
		time.Sleep(100 * time.Millisecond)
	}()

	wg.Wait()
	fmt.Println("Seguro: sin deadlock.")
}

func main() {
	fmt.Println("=== Elige una sección para ejecutar ===")
	// comenta/activa la versión que desees probar

	// deadlock()      // <- provocará interbloqueo
	seguroOrdenado()   // <- versión segura
}
