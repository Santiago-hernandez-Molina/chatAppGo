package email

import (
	"fmt"
	"net/smtp"

	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/models"
	"github.com/Santiago-hernandez-Molina/chatAppBackend/internal/domain/ports"
)

type EmailSender struct {
	emailServer string
	password    string
	provider    string
}

func (sender *EmailSender) SendRegisterConfirm(user *models.User, code int) error {
	auth := smtp.PlainAuth("", sender.emailServer, sender.password, sender.provider)
	to := []string{user.Email}
	body := fmt.Sprintf("From: '%s ✉️' <%s>\nSubject:%s\n%s%d",
		"ChatApp",
		sender.emailServer,
		"ChatApp Register Code",
        "This is your activation code, don't share with anyone!! Expires in 10 minutes: ",
        code,
	)
	msg := []byte(body)

	err := smtp.SendMail(sender.provider+":587", auth, sender.emailServer, to, msg)
	if err != nil {
		return err
	}
	return nil
}

var _ ports.EmailSender = (*EmailSender)(nil)

func NewEmailSender(emailSender string, password string, provider string) *EmailSender {
	return &EmailSender{
		emailServer: emailSender,
		password:    password,
		provider:    provider,
	}
}
