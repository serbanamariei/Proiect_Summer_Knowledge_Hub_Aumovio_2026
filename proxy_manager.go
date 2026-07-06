package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
)

type ProxyManager struct {
	pool []ProxyData
	mu   sync.Mutex
}

var (
	instance *ProxyManager
	once     sync.Once
	initErr  error
)

func GetProxyManager(db *sql.DB) (*ProxyManager, error) {
	once.Do(func() {
		proxyManager := &ProxyManager{
			pool: []ProxyData{},
		}

		rows, err := db.Query(`SELECT type, address, username, password FROM proxyuri`)
		if err != nil {
			initErr = fmt.Errorf("eroare la interogarea proxy-urilor: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var p ProxyData
			if err := rows.Scan(&p.Tip, &p.Adresa, &p.Nume, &p.Parola); err != nil {
				log.Println("eroare la citirea proxy-ului din baza de date:", err)
				continue
			}
			proxyManager.pool = append(proxyManager.pool, p)
		}

		if len(proxyManager.pool) == 0 {
			initErr = fmt.Errorf("baza de date nu are proxy-uri")
			return
		}

		log.Printf("in ProxyManager sunt %d proxy-uri\n", len(proxyManager.pool))
		instance = proxyManager
	})

	if initErr != nil {
		return nil, initErr
	}

	return instance, nil
}

func (pm *ProxyManager) GetRandomProxy() ProxyData {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	randomIndex := rand.Intn(len(pm.pool))
	return pm.pool[randomIndex]
}

type ProxyFilter func(p ProxyData) bool

func FiltruSanatateStandard(p ProxyData) bool {
	return !p.EBanat && p.NrEsuari <= 5
}

func FiltruStrict(p ProxyData) bool {
	return !p.EBanat && p.NrEsuari == 0 && p.NrSuccese > 5
}

func (pm *ProxyManager) GetBestProxy(filter ProxyFilter) ProxyData {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var proxyUriSanatoase []ProxyData

	for _, p := range pm.pool {
		if filter(p) {
			proxyUriSanatoase = append(proxyUriSanatoase, p)
		}
	}

	if len(proxyUriSanatoase) == 0 {
		randomIndex := rand.Intn(len(pm.pool))
		return pm.pool[randomIndex]
	}

	randomIndex := rand.Intn(len(proxyUriSanatoase))
	return proxyUriSanatoase[randomIndex]
}
