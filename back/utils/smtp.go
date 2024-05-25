package utils

import (
    "fmt"
    "net/smtp"
)

func SendEmail(to string, subject string, body string) error {
    smtpHost := "smtp.example.com"
    smtpPort := "587"

    from := "your-email@example.com"
    password := "your-email-password"

    msg := "From: " + from + "\n" +
        "To: " + to + "\n" +
        "Subject: " + subject + "\n\n" +
        body

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
    if err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }
    return nil
}
