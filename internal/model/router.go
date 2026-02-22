package model

// RouterConfig represents a configured MikroTik router session.
type RouterConfig struct {
	SessionName string `json:"session_name"`
	Host        string `json:"host"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	HotspotName string `json:"hotspot_name"`
	DNSName     string `json:"dns_name"`
	Currency    string `json:"currency"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	InfoLP      string `json:"info_lp"`
	IdleTimeout string `json:"idle_timeout"`
	ReportFlag  string `json:"report_flag"`
	Token       string `json:"token"`
}
