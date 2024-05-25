package utils

import (
    "fmt"
    "net/smtp"
)

// SendEmail envía un correo electrónico utilizando el servidor SMTP especificado.
func SendEmail(to string, subject string, body string) error {
    // Dirección del servidor SMTP
    smtpHost := "smtp.example.com"
    smtpPort := "587"

    // Credenciales de autenticación SMTP
    from := "your-email@example.com"
    password := "your-email-password"

    // Crear el mensaje
    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    // Dirección de envío
    auth := smtp.PlainAuth("", from, password, smtpHost)

    // Enviar el correo electrónico
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }
    return nil
}
