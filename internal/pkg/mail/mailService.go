package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

type PostService struct {
	dialer *gomail.Dialer
	from   string
}

func NewMailService(host string, port int, username, password, from string) *PostService {
	return &PostService{
		dialer: gomail.NewDialer(host, port, username, password),
		from:   from,
	}
}

func (s *PostService) SendConfirmationEmail(email, token string) error {

	domain := os.Getenv("DOMAIN")

	subject := "Подтверждение регистрации на Lootor"
	body := fmt.Sprintf("Подтвердите ваш аккаунт, перейдя по ссылке: https://%s/auth?confirmationToken=%s", domain, token)
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return s.dialer.DialAndSend(m)
}

func (s *PostService) SendResetPasswordEmail(email, token string) error {
	domain := os.Getenv("DOMAIN")

	subject := "Восстановление пароля на Lootor"
	body := fmt.Sprintf("Восстановите пароль, перейдя по ссылке: https://%s/auth/reset?resetToken=%s \nЕсли вы не пытались восстановить пароль, то проигнорируйте это письмо", domain, token)
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	return s.dialer.DialAndSend(m)
}
