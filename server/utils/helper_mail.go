package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// SendResetEmail gửi email reset qua Gmail SMTP (yêu cầu App Password)
func SendResetEmail(toEmail, resetLink string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	if from == "" || pass == "" {
		return fmt.Errorf("SMTP_EMAIL or SMTP_PASSWORD not set")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	// Email headers and body (plain text + simple HTML)
	subject := "Reset your password"
	bodyText := fmt.Sprintf("Reset your password using this link:\n\n%s\n\nIf you didn't request this, ignore.", resetLink)
	bodyHTML := fmt.Sprintf("<p>Reset your password using this link:</p><p><a href=\"%s\">Reset password</a></p>", resetLink)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: multipart/alternative; boundary=BOUNDARY\r\n")
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("\r\n--BOUNDARY\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	msg.WriteString(bodyText + "\r\n")
	msg.WriteString("\r\n--BOUNDARY\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	msg.WriteString(bodyHTML + "\r\n")
	msg.WriteString("\r\n--BOUNDARY--")

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	if err := smtp.SendMail(addr, auth, from, []string{toEmail}, []byte(msg.String())); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
