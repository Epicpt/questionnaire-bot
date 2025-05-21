package smtp

import (
	"fmt"
	"questionnaire-bot/internal/config"
	"questionnaire-bot/internal/entity"
	"strings"

	"gopkg.in/gomail.v2"
)

const Subject = "Новый лид с опросника Telegram Bot"

type EmailClient interface {
	Send(*entity.Email) error
}

type MailClient struct {
	recipient string
	from      string
	Dialer    *gomail.Dialer
}

func New(c config.SMTP) *MailClient {
	return &MailClient{
		recipient: c.Client,
		from:      c.Sender,
		Dialer:    gomail.NewDialer(c.Host, c.Port, c.Username, c.Password),
	}
}

func (c *MailClient) Send(email *entity.Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", c.from)
	m.SetHeader("To", c.recipient)
	m.SetHeader("Subject", Subject)
	m.SetBody("text/plain", email.Body)

	htmlBody := strings.ReplaceAll(email.Body, "\n", "<br>")
	m.AddAlternative("text/html", htmlBody)

	if err := c.Dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
