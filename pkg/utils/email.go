package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/TFX0019/api-go-gds/pkg/config"
	"github.com/resend/resend-go/v3"
)

func GenerateSixDigitCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(999999)
	return fmt.Sprintf("%06d", code)
}

func SendVerificationEmail(email, token string) {
	smtpHost := config.GetEnv("SMTP_HOST", "")
	smtpPort := config.GetEnv("SMTP_PORT", "")
	from := config.GetEnv("SMTP_EMAIL", "")
	password := config.GetEnv("SMTP_PASSWORD", "")

	to := []string{email}
	subject := "Verify your account"

	baseURL := config.GetEnv("API_URL", "http://localhost:3000")
	code := fmt.Sprintf("BODY: Click here to verify: %s/api/auth/verify?token=%s", baseURL, token)

	message := []byte("To: " + to[0] + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		code)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if err != nil {
		log.Printf("Error al enviar el email: %v", err)
		if strings.Contains(err.Error(), "535") {
			log.Println("Sugerencia: El error 535 suele indicar un fallo de autenticación. ¡Asegúrate de usar la Contraseña de Aplicación de 16 caracteres!")
		}
		return
	}

	// In a real application, use an SMTP server here.
	log.Printf("----------------------------------------------------------------")
	log.Printf("EMAIL SENT TO: %s", email)
	log.Printf("SUBJECT: Verify your account")
	log.Printf("BODY: Click here to verify: /api/auth/verify?token=%s", token)
	log.Printf("----------------------------------------------------------------")
}

func SendRecoveryEmail(email, code string) {
	smtpHost := config.GetEnv("SMTP_HOST", "smtp.gmail.com")
	smtpPort := config.GetEnv("SMTP_PORT", "587")
	from := config.GetEnv("SMTP_EMAIL", "1400kss@gmail.com")
	password := config.GetEnv("SMTP_PASSWORD", "nknf xqed dszx swoe")

	to := []string{email}
	subject := "Verify your account"

	codeBody := fmt.Sprintf("BODY: Your recovery code is: %s", code)

	message := []byte("To: " + to[0] + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		codeBody)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if err != nil {
		log.Printf("Error al enviar el email: %v", err)
		if strings.Contains(err.Error(), "535") {
			log.Println("Sugerencia: El error 535 suele indicar un fallo de autenticación. ¡Asegúrate de usar la Contraseña de Aplicación de 16 caracteres!")
		}
		return
	}

	// In a real application, use an SMTP server here.
	log.Printf("----------------------------------------------------------------")
	log.Printf("EMAIL SENT TO: %s", email)
	log.Printf("SUBJECT: Password Recovery")
	log.Printf("BODY: Your recovery code is: %s", code)
	log.Printf("----------------------------------------------------------------")
}

func SendTestEmail(toEmail string) error {
	apiKey := config.GetEnv("RESEND_API_KEY", "")
	if apiKey == "" {
		return fmt.Errorf("RESEND_API_KEY is not set")
	}

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <noreply@patronesparacostura.com>",
		To:      []string{toEmail},
		Html:    "<strong>hello world</strong>",
		Subject: "Hello from Golang",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	log.Printf("Test Email Sent. ID: %s", sent.Id)
	return nil
}
