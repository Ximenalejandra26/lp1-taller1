package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Objetivo: Lectores–Escritores con RWMutex sobre un mapa compartido.
// Varios lectores pueden leer en paralelo; los escritores tienen exclusión mutua.
// completa los pasos y observa la diferencia entre Mutex y RWMutex.

type baseDatos struct {
	mu sync.RWMutex // cambia a sync.Mutex para comparar comportamiento
	m  map[string]int
}

func (db *baseDatos) leer(clave string) (int, bool) {
	//  usar RLock/RUnlock (o Lock/Unlock si usas Mutex)
	db.mu.RLock()		
	defer db.mu.RUnlock()

	v, ok := db.m[clave]
	return v, ok
}

func (db *baseDatos) escribir(clave string, valor int) {
	//  usar Lock/Unlock para escritura
	db.mu.Lock()
	defer db.mu.Unlock()

	db.m[clave] = valor
}

func lector(id int, db *baseDatos, claves []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		clave := claves[rand.Intn(len(claves))]
		if v, ok := db.leer(clave); ok {
			fmt.Printf("[lector %d] %s=%d\n", id, clave, v)
		} else {
			fmt.Printf("[lector %d] %s no existe\n", id, clave)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func escritor(id int, db *baseDatos, claves []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		clave := claves[rand.Intn(len(claves))]
		v := rand.Intn(1000)
		db.escribir(clave, v)
		fmt.Printf("[escritor %d] set %s=%d\n", id, clave, v)
		time.Sleep(120 * time.Millisecond)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	db := &baseDatos{m: make(map[string]int)}
	claves := []string{"a", "b", "c", "d", "e"}

	// precarga
	for _, k := range claves {
		db.m[k] = rand.Intn(100)
	}

	var wg sync.WaitGroup

	// experimenta con # de lectores y escritores
	nLectores := 5
	nEscritores := 2

	wg.Add(nLectores + nEscritores)
	for i := 1; i <= nLectores; i++ {
		go lector(i, db, claves, &wg)
	}
	for j := 1; j <= nEscritores; j++ {
		go escritor(j, db, claves, &wg)
	}

	wg.Wait()
	fmt.Println("Fin Lectores–Escritores.")
}
