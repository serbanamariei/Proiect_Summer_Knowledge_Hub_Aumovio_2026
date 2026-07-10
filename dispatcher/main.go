package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"comun"

	_ "github.com/lib/pq"
)

var db *sql.DB
var pm *ProxyManager

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	pm, err = NewProxyManager(db)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/dispatch", handleDispatch)
	http.ListenAndServe("0.0.0.0:8082", nil)
}

func handleDispatch(w http.ResponseWriter, r *http.Request) {
	var cerere comun.CerereInitala
	json.NewDecoder(r.Body).Decode(&cerere)

	proxyAles := pm.GetBestProxy()

	cerereCrawler := comun.CerereCrawler{
		ChatID: cerere.ChatID,
		URL:    cerere.URL,
		Proxy:  proxyAles,
	}

	reqBody, _ := json.Marshal(cerereCrawler)

	resp, err := http.Post("http://crawler_aplicatie:8081/proceseaza", "application/json", bytes.NewBuffer(reqBody))

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"Mesaj": "Eroare: Crawler offline."})
		return
	}
	defer resp.Body.Close()

	var raspuns comun.RaspunsCrawler
	json.NewDecoder(resp.Body).Decode(&raspuns)

	linkuriJSON, err := json.Marshal(raspuns.Linkuri)
	if err == nil {
		db.Exec("INSERT INTO siteuri_parsate (url, linkuri) VALUES ($1, $2)", cerere.URL, linkuriJSON)
	}

	json.NewEncoder(w).Encode(raspuns)
}
