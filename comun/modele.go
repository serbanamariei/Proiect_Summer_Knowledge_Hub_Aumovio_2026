package comun

type ProxyData struct {
	Tip       string `json:"type"`
	Adresa    string `json:"address"`
	Nume      string `json:"username"`
	Parola    string `json:"password"`
	NrSuccese int    `json:"success_count"`
	NrEsuari  int    `json:"fail_count"`
	EBanat    bool   `json:"is_banned"`
}

type IPResponse struct {
	IP string `json:"ip"`
}

type CerereCrawler struct {
	ChatID int64     `json:"chat_id"`
	URL    string    `json:"url"`
	Proxy  ProxyData `json:"proxy"`
}

type RaspunsCrawler struct {
	Mesaj   string   `json:"mesaj"`
	Linkuri []string `json:"linkuri"`
}

type CerereInitala struct {
	ChatID int64  `json:"chat_id"`
	URL    string `json:"url"`
}
