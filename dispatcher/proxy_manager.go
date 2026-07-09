package main

import (
	"comun"
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
)

type ProxyManager struct {
	pool      []comun.ProxyData
	sanatoase []comun.ProxyData
	mu        sync.Mutex
}

var (
	instance *ProxyManager
	once     sync.Once
	initErr  error
)

func NewProxyManager(db *sql.DB) (*ProxyManager, error) {
	once.Do(func() {
		pm := &ProxyManager{
			pool:      []comun.ProxyData{},
			sanatoase: []comun.ProxyData{},
		}

		rows, err := db.Query(`SELECT tip, adresa, nume, parola, nr_succese, nr_esuari, e_banat FROM proxyuri`)
		if err != nil {
			initErr = fmt.Errorf("eroare DB: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var proxyCurent comun.ProxyData
			if err := rows.Scan(&proxyCurent.Tip, &proxyCurent.Adresa, &proxyCurent.Nume, &proxyCurent.Parola, &proxyCurent.NrSuccese, &proxyCurent.NrEsuari, &proxyCurent.EBanat); err != nil {
				continue
			}
			pm.pool = append(pm.pool, proxyCurent)

			if !proxyCurent.EBanat && proxyCurent.NrEsuari <= 5 {
				pm.sanatoase = append(pm.sanatoase, proxyCurent)
			}
		}
		instance = pm
	})

	if initErr != nil {
		return nil, initErr
	}
	return instance, nil
}

func (pm *ProxyManager) GetBestProxy() comun.ProxyData {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if len(pm.sanatoase) > 0 {
		return pm.sanatoase[rand.Intn(len(pm.sanatoase))]
	}

	if len(pm.pool) > 0 {
		return pm.pool[rand.Intn(len(pm.pool))]
	}

	return comun.ProxyData{}
}
