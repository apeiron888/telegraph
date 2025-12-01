package auth

import (
	"fmt"
	"net/smtp"
)

type EmailSender interface {
	Send(to, subject, body string) error
}

type SMTPSender struct {
	from     string
	password string
	host     string
	port     string
}

func NewSMTPSender(from, password, host, port string) *SMTPSender {
	return &SMTPSender{from, password, host, port}
}

func (s *SMTPSender) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.from, s.password, s.host)

	msg := "From: " + s.from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}
