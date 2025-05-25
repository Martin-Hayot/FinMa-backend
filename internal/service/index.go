package service

type Services struct {
	Auth        AuthService
	User        UserService
	GoCardless  GoCardlessService
	BankAccount BankAccountService
	Validator   ValidatorService
}
