package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

func GenerateSixDigitCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(999999)
	return fmt.Sprintf("%06d", code)
}

func SendVerificationEmail(email, token string) {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	from := "1400kss@gmail.com"
	password := "nknf xqed dszx swoe"

	to := []string{email}
	subject := "Verify your account"

	code := fmt.Sprintf("BODY: Click here to verify: http://localhost:3000/api/auth/verify?token=%s", token)

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
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	from := "1400kss@gmail.com"
	password := "nknf xqed dszx swoe"

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
