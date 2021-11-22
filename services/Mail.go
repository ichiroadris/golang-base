package services

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"os"
)

type EmailObject struct {
	To      string
	Body    string
	Subject string
}

func SendMail(subject string, body string, to string, html string, name string) bool {

	from := mail.NewEmail("Just Open it", os.Getenv("SENDGRID_FROM_EMAIL"))

	_to := mail.NewEmail(name, to)

	plainTextContext := body

	htmlContent := html
	message := mail.NewSingleEmail(from, subject, _to, plainTextContext, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		return false
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return true
	}
}
