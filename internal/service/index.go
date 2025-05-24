package service

type Services struct {
	Auth        AuthService
	User        UserService
	BankAccount BankAccountService
	PlaidItem   PlaidItemService
	Validator   ValidatorService
}
