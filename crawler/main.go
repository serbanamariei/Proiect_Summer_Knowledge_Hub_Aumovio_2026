package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"comun"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/proceseaza", handleProceseaza)
	http.ListenAndServe("0.0.0.0:8081", nil)
}

type IProxy interface {
	ConstruiesteClient() (*http.Client, error)
}

type SocksProxy struct {
	Tip    string
	Adresa string
	Nume   string
	Parola string
}

func (s *SocksProxy) ConstruiesteClient() (*http.Client, error) {
	var proxyString string
	if s.Nume != "" || s.Parola != "" {
		proxyString = fmt.Sprintf("%s://%s:%s@%s", s.Tip, s.Nume, s.Parola, s.Adresa)
	} else {
		proxyString = fmt.Sprintf("%s://%s", s.Tip, s.Adresa)
	}

	proxyURL, err := url.Parse(proxyString)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
	}, nil
}

type NoProxy struct{}

func (n *NoProxy) ConstruiesteClient() (*http.Client, error) {
	return &http.Client{}, nil
}

func CreazaProxyConcret(date comun.ProxyData) (IProxy, error) {
	if date.Adresa == "" {
		return &NoProxy{}, nil
	}
	switch date.Tip {
	case "socks4", "socks5":
		return &SocksProxy{Tip: date.Tip, Adresa: date.Adresa, Nume: date.Nume, Parola: date.Parola}, nil
	default:
		return nil, fmt.Errorf("tip de proxy nesuportat momentan: %s", date.Tip)
	}
}

func handleProceseaza(w http.ResponseWriter, r *http.Request) {
	var cerere comun.CerereCrawler
	json.NewDecoder(r.Body).Decode(&cerere)

	w.Header().Set("Content-Type", "application/json")

	proxyConcret, err := CreazaProxyConcret(cerere.Proxy)
	if err != nil {
		json.NewEncoder(w).Encode(comun.RaspunsCrawler{Mesaj: "Eroare instantiere proxy: " + err.Error()})
		return
	}

	client, err := proxyConcret.ConstruiesteClient()
	if err != nil {
		json.NewEncoder(w).Encode(comun.RaspunsCrawler{Mesaj: "Eroare construire client HTTP: " + err.Error()})
		return
	}

	htmlBrut, err := ExtrageHTML(client, cerere.URL)
	if err != nil {
		json.NewEncoder(w).Encode(comun.RaspunsCrawler{Mesaj: "Eroare Scraping: " + err.Error()})
		return
	}

	listaLinkuri := ExtrageLinkuri(htmlBrut)

	for i := 0; i < len(listaLinkuri) && i < 10; i++ {
		db.Exec(`INSERT INTO rezultate_extragere (link_sursa, link_gasit) VALUES ($1, $2)`, cerere.URL, listaLinkuri[i])
	}

	mesaj := fmt.Sprintf("Scraping efectuat. Am gasit %d link-uri.", len(listaLinkuri))
	if cerere.Proxy.Adresa != "" {
		mesaj = fmt.Sprintf("Scraping efectuat prin %s. Am gasit %d link-uri.", cerere.Proxy.Adresa, len(listaLinkuri))
	}

	raspuns := comun.RaspunsCrawler{
		Mesaj:   mesaj,
		Linkuri: listaLinkuri,
	}

	json.NewEncoder(w).Encode(raspuns)
}

func ExtrageHTML(client *http.Client, link string) (string, error) {
	resp, err := client.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func ExtrageLinkuri(html string) []string {
	re := regexp.MustCompile(`href="(http[s]?://[^"]+)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	linkuriUnice := make(map[string]bool)
	var listaFinala []string

	for _, match := range matches {
		if len(match) > 1 {
			link := match[1]
			if !linkuriUnice[link] {
				linkuriUnice[link] = true
				listaFinala = append(listaFinala, link)
			}
		}
	}

	return listaFinala
}
