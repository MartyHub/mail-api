package smtp

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	AuthPlain   = "PLAIN"
	AuthCRAMMD5 = "CRAMMD5"
)

type Service interface {
	Send(mail Mail) error
}

func NewService(cfg Config) (Service, error) {
	result := service{
		addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		cfg:  cfg,
	}

	switch cfg.AuthType {
	case AuthCRAMMD5:
		result.auth = smtp.CRAMMD5Auth(cfg.AuthUsername, cfg.AuthPassword)
	case AuthPlain:
		result.auth = smtp.PlainAuth("", cfg.AuthUsername, cfg.AuthPassword, cfg.Host)
	}

	return result, nil
}

type service struct {
	addr string
	auth smtp.Auth
	cfg  Config
}

type Mail struct {
	From    string
	To      []string
	Cc      []string
	Subject string
	Body    string
	HTML    bool
}

func (m Mail) msg() []byte {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.To, ",")))
	sb.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(m.Cc, ",")))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", m.Subject))

	if m.HTML {
		sb.WriteString("MIME-Version: 1.0\r\n")
		sb.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	}

	sb.WriteString("\r\n")
	sb.WriteString(m.Body)

	return []byte(sb.String())
}

func (s service) Send(m Mail) error {
	return smtp.SendMail(
		s.addr,
		s.auth,
		s.cfg.From,
		m.To,
		m.msg(),
	)
}
