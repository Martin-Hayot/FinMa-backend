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

type LinkAccountRequest struct {
	InstitutionID string `json:"institution_id"`
}

type LinkAccountResponse struct {
	Link string `json:"link"` // URL to redirect user to for linking account
}

// Internal GoCardless API DTO (for communicating with GoCardless)
type GoCardlessCreateRequisitionRequest struct {
	InstitutionID string `json:"institution_id"`
	RedirectURL   string `json:"redirect"`
	Reference     string `json:"reference"` // GoCardless expects string
}

// Internal GoCardless API DTO (for communicating with GoCardless)
type GoCardlessCreateRequisitionResponse struct {
	ID                string   `json:"id"`
	Created           string   `json:"created"`
	RedirectURL       string   `json:"redirect"`
	Status            string   `json:"status"`
	InstitutionID     string   `json:"institution_id"`
	Agreement         string   `json:"agreement"`
	Reference         string   `json:"reference"`
	Accounts          []string `json:"accounts"`
	UserLanguage      string   `json:"user_language"`
	Link              string   `json:"link"`
	SSN               string   `json:"ssn"`
	AccountSelection  bool     `json:"account_selection"`
	RedirectImmediate bool     `json:"redirect_immediate"`
}

type GoCardlessGetRequisitionResponse struct {
	ID                string   `json:"id"`
	Created           string   `json:"created"`
	RedirectURL       string   `json:"redirect"`
	Status            string   `json:"status"`
	InstitutionID     string   `json:"institution_id"`
	Agreement         string   `json:"agreement"`
	Reference         string   `json:"reference"`
	Accounts          []string `json:"accounts"`
	UserLanguage      string   `json:"user_language"`
	Link              string   `json:"link"`
	SSN               string   `json:"ssn"`
	AccountSelection  bool     `json:"account_selection"`
	RedirectImmediate bool     `json:"redirect_immediate"`
}

type GoCardlessGetRequisitionRequest struct {
	RequisitionID string `json:"requisition_id"`
}

type GoCardlessUpdateRequisitionRequest struct {
	RequisitionReference string `json:"requisition_reference"` // Reference to update the requisition
}
type GoCardlessUpdateRequisitionResponse struct {
	Status        string `json:"status"`
	InstitutionID string `json:"institution_id"`
	Reference     string `json:"reference"`
}
