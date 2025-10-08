package utils

import (
	"auth-service/config"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

func SendEmailOTP(to, otp, lang string) error {
	mailer := config.GetConfig().MAILER
	// Load template HTML
	tmpl, err := template.ParseFiles(fmt.Sprintf("templates/otp_email_%s.html", lang))
	if err != nil {
		return fmt.Errorf("cannot load template: %w", err)
	}

	var body bytes.Buffer
	data := struct{ OTP string }{OTP: otp}
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("cannot execute template: %w", err)
	}

	msg := "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mailer.From)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += "Subject: Mã xác thực OTP của bạn\r\n\r\n"
	msg += body.String()

	auth := smtp.PlainAuth("", mailer.From, mailer.Password, mailer.Host)

	addr := fmt.Sprintf("%s:%s", mailer.Host, mailer.Port)
	if err := smtp.SendMail(addr, auth, mailer.From, []string{to}, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("OTP email sent to: ", to)
	return nil
}
