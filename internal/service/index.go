package service

type Services struct {
	Auth        AuthService
	User        UserService
	GoCardless  GclService
	BankAccount BankAccountService
	Validator   ValidatorService
}
