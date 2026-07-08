package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

func ExtrageHTML(link string, proxy ProxyData) (string, error) {
	var proxyString string
	if proxy.Nume != "" && proxy.Parola != "" {
		proxyString = fmt.Sprintf("%s://%s:%s@%s", proxy.Tip, proxy.Nume, proxy.Parola, proxy.Adresa)
	} else {
		proxyString = fmt.Sprintf("%s://%s", proxy.Tip, proxy.Adresa)
	}

	proxyURL, err := url.Parse(proxyString)
	if err != nil {
		return "", fmt.Errorf("eroare URL proxy: %w", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ro;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status cod eroare: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func ExtrageLinkuri(codHTML string) []string {
	regex := regexp.MustCompile(`href=["']([^"']+)["']`)
	potriviri := regex.FindAllStringSubmatch(codHTML, -1)

	linkuriUnice := make(map[string]bool)
	var listaFinala []string

	for _, potrivire := range potriviri {
		if len(potrivire) > 1 {
			link := potrivire[1]
			if !linkuriUnice[link] {
				linkuriUnice[link] = true
				listaFinala = append(listaFinala, link)
			}
		}
	}

	return listaFinala
}
