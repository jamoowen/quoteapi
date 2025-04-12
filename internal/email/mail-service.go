package email

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/jamoowen/quoteapi/internal/utils"
)

type MailService interface {
	Send(to, subject, message string) error
}

type smtpMailer struct {
	from string
	addr string
	auth smtp.Auth
}

func NewSmtpMailer(from, password, smtpHost, smtpPort string) (*smtpMailer, error) {
	// hostname is used by PlainAuth to validate the TLS certificate.
	if utils.LooksLikeEmail(from) == false {
		return nil, fmt.Errorf("from email looks bad: %v", from)
	}
	if password == "" || smtpHost == "" || smtpPort == "" {
		return nil, fmt.Errorf("from(%v), password(%v), smtpHost(%v), smtpPort(%v) are all required", from, password, smtpHost, smtpPort)
	}
	auth := smtp.PlainAuth("", from, password, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	return &smtpMailer{
		from: from,
		addr: addr,
		auth: auth,
	}, nil
}

func (s *smtpMailer) Send(to, subject, message string) error {
	htmlBody := fmt.Sprintf(`
        <html>
        <head>
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
        </head>
        <body style="font-family: Arial, sans-serif; line-height: 1.6; margin: 0; padding: 20px;">
            <div style="max-width: 600px; margin: 0 auto; background: #ffffff;">
                <div style="padding: 20px;">
                    <p>%s</p>
                </div>
                <div style="padding: 20px; font-size: 12px; color: #666; border-top: 1px solid #eee;">
                    <p>This is an automated message. Please do not reply.</p>
                    <p>If you didn't request this email, please ignore it.</p>
                </div>
            </div>
        </body>
        </html>
    `, message)

	// Construct email with headers
	emailMsg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"X-Mailer: GoMailer\r\n"+
		"Message-ID: <%s@%s>\r\n"+ // Add unique message ID
		"\r\n"+
		"%s",
		s.from,
		to,
		subject,
		generateMessageID(),           // implement this function
		strings.Split(s.from, "@")[1], // domain part of email
		htmlBody))

	// Send email
	err := smtp.SendMail(
		s.addr,
		s.auth,
		s.from,
		[]string{to},
		emailMsg,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// Generate a unique message ID
func generateMessageID() string {
	return fmt.Sprintf("%d.%d", time.Now().UnixNano(), rand.Int63())
}
