package mail

import (
	"Market_backend/internal/config"
	"crypto/tls"
	"fmt"
	"net/smtp"
)

type MailService struct {
	Host     string
	Port     string
	Email    string
	Password string
	From     string
}

func NewMailService() *MailService {
	return &MailService{
		Host:     config.SMTPHost,
		Port:     config.SMTPPort,
		Email:    config.SMTPEmail,
		Password: config.SMTPPassword,
		From:     config.SMTPEmail, // если from = тот же email
	}
}

func (m *MailService) SendEmail(to, subject, body string) error {
	// Адрес SMTP сервера
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	// Создаем auth
	auth := smtp.PlainAuth("", m.Email, m.Password, m.Host)

	// Создаем TCP соединение
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Quit()

	// Обновляем соединение в TLS (STARTTLS)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, // только для dev
		ServerName:         m.Host,
	}
	if err = c.StartTLS(tlsconfig); err != nil {
		return err
	}

	// Авторизация
	if err = c.Auth(auth); err != nil {
		return err
	}

	// От кого
	if err = c.Mail(m.From); err != nil {
		return err
	}

	// Кому
	if err = c.Rcpt(to); err != nil {
		return err
	}

	// Отправка письма
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	msg := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)
	_, err = wc.Write([]byte(msg))
	return err
}
