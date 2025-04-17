package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type MailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailService(host string, port int, username, password, from string) *MailService {
	return &MailService{
		dialer: gomail.NewDialer(host, port, username, password),
		from:   from,
	}
}

func (s *MailService) SendConfirmationEmail(email, token string) error {

	subject := "Подтверждение регистрации на Lootor"
	body := fmt.Sprintf("Подтвердите ваш аккаунт, перейдя по ссылке: https://dev.lootor.me/auth?confirmationToken=%s", token)
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return s.dialer.DialAndSend(m)
}
