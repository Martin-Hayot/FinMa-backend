package utils

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

// Initialisation de la clé API et du client une seule fois
var (
    resendAPIKey = os.Getenv("RESEND_API_KEY")
    resendClient = resend.NewClient(resendAPIKey)
)

func SendMail(email string, subject string, body string) error {
    params := &resend.SendEmailRequest{
        From:    "FinMa <noreply@gaetanleplae.com>",
        To:      []string{email},
        Subject: subject,
        Html:    body,
    }
	fmt.Println("clé API", resendAPIKey)

    sent, err := resendClient.Emails.Send(params)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }
    fmt.Println(sent.Id)
    return nil
}