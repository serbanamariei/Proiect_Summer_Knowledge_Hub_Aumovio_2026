package main

type TelegramUpdate struct {
	Message TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	Chat TelegramChat `json:"chat"`
	Text string       `json:"text"`
}

type TelegramChat struct {
	ID int `json:"id"`
}

type SendMessageReq struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type ProxyData struct {
	Tip    string `json:"type"`
	Adresa string `json:"address"`
	Nume   string `json:"username"`
	Parola string `json:"password"`

	NrSuccese int  `json:"success_count"`
	NrEsuari  int  `json:"fail_count"`
	EBanat    bool `json:"is_banned"`
}

type CrawlerConfig struct {
	TargetURL string
	Proxy     ProxyData
}

type CrawlerResult struct {
	ProxyIP string
	Linkuri []string
	Eroare  error
}

type IPResponse struct {
	IP string `json:"ip"`
}
