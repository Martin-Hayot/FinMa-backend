package dto

type Institution struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	BIC       string   `json:"bic"`
	Countries []string `json:"countries"`
	Logo      string   `json:"logo"`
}

type TokenRequest struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

type TokenResponse struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"access_expires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refresh_expires"`
}
