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
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Email, m.Password, m.Host)

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Quit()

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true, // только для dev
		ServerName:         m.Host,
	}
	if err = c.StartTLS(tlsconfig); err != nil {
		return err
	}

	if err = c.Auth(auth); err != nil {
		return err
	}

	if err = c.Mail(m.From); err != nil {
		return err
	}

	if err = c.Rcpt(to); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	// Здесь добавляем заголовок Content-Type: text/html
	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n\r\n%s",
		m.From, to, subject, body,
	)

	_, err = wc.Write([]byte(msg))
	return err
}
