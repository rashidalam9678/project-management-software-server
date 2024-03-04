package helpers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendInviteMail(email string, roleID string, token string) (res interface{}, err error) {
	godotenv.Load()
	sendgrid_key := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("Taskify", "gj9678@myamu.ac.in")
	subject := "Project Member Invitation"
	to := mail.NewEmail("Example User", email)
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := fmt.Sprintf(`<h1>Hello from Taskify</h1>, You have been invited to join a project. Click <a href="http://localhost:3000/invite?token=%s&role=%s">here</a> to accept the invitation.`, token,roleID)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendgrid_key)
	_, err = client.Send(message)
	if err != nil {
		log.Println(err)
		return nil, err
	} 
	return nil, nil
}
