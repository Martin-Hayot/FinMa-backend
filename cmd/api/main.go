package main

import (
	"FinMa/internal/server"
	"FinMa/utils"
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

type Payload struct {
	UserID   int
	Username string
	Email    string
}

func main() {

	server := server.New()

	server.RegisterFiberRoutes()

	// Generate an access token
	payload := Payload{
		UserID:   1,
		Username: "john_doe",
		Email:    "john.doe@example.com",
	}

	token, err := utils.GenerateAccessToken(payload)
	if err != nil {
		panic(fmt.Sprintf("cannot generate access token: %s", err))
	}

	fmt.Printf("Access Token: %s\n", token)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err = server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
