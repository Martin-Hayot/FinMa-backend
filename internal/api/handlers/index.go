package handlers

type Handlers struct {
	Auth        AuthHandler
	User        UserHandler
	GoCardless  GclHandler
	BankAccount BankAccountHandler
}
