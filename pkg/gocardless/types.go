package gocardless

type Institution struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	BIC       string   `json:"bic"`
	Countries []string `json:"countries"`
	Logo      string   `json:"logo"`
}

type TokenRequest struct {
	SecretID  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type TokenResponse struct {
	Access         string `json:"access"`
	AccessExpires  int    `json:"accessExpires"`
	Refresh        string `json:"refresh"`
	RefreshExpires int    `json:"refreshExpires"`
}

type LinkAccountRequest struct {
	InstitutionID string `json:"institutionId"`
}

type LinkAccountResponse struct {
	Link string `json:"link"` // URL to redirect user to for linking account
}

// GoCardlessCreateRequisitionRequest is the request body for creating a requisition
type GoCardlessCreateRequisitionRequest struct {
	InstitutionID string `json:"institutionId"`
	Redirect      string `json:"redirect"`
	Reference     string `json:"reference"`
}

// AccountDetails represents the details of a bank account
type AccountDetails struct {
	Account struct {
		ResourceID      string `json:"resourceId"`
		IBAN            string `json:"iban"`
		Currency        string `json:"currency"`
		OwnerName       string `json:"ownerName"`
		Name            string `json:"name"`
		Product         string `json:"product"`
		Status          string `json:"status"`
		InstitutionName string `json:"institutionName"`
	} `json:"account"`
}

// AccountBalance represents a single balance entry
type AccountBalance struct {
	BalanceAmount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"balanceAmount"`
	BalanceType      string `json:"balanceType"`
	ReferenceDate    string `json:"referenceDate"`
	LastChange       string `json:"lastChangeDateTime"`
	LastCommittedTxn string `json:"lastCommittedTransaction"`
}

// AccountBalances represents the balances of a bank account
type AccountBalances struct {
	Balances []AccountBalance `json:"balances"`
}

// Transaction represents a single transaction
type Transaction struct {
	TransactionID     string `json:"transactionId"`
	BookingDate       string `json:"bookingDate"`
	ValueDate         string `json:"valueDate"`
	BookingDateTime   string `json:"bookingDateTime"`
	ValueDateTime     string `json:"valueDateTime"`
	TransactionAmount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"transactionAmount"`
	CreditorName          string `json:"creditorName"`
	DebtorName            string `json:"debtorName"`
	RemittanceInformation string `json:"remittanceInformationUnstructured"`
}

// AccountTransactions represents the transactions of a bank account
type AccountTransactions struct {
	Transactions struct {
		Booked  []Transaction `json:"booked"`
		Pending []Transaction `json:"pending"`
	} `json:"transactions"`
}

// Internal GoCardless API DTO (for communicating with GoCardless)
type GoCardlessCreateRequisitionResponse struct {
	ID                string   `json:"id"`
	Created           string   `json:"created"`
	Redirect          string   `json:"redirect"`
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
	Redirect          string   `json:"redirect"`
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
