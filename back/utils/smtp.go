package utils

import (
    "crypto/tls"
    "fmt"
    "log"
    "net/smtp"
    "net/url"
    "os"

    "github.com/joho/godotenv"
)

func LoadEnv() error {
    return godotenv.Load()
}

func SendEmail(to string, subject string, body string) error {
    err := LoadEnv()
    if err != nil {
        return fmt.Errorf("failed to load .env file: %v", err)
    }

    mailerURL, err := url.Parse(os.Getenv("MAILER"))
    if err != nil {
        return fmt.Errorf("failed to parse MAILER URL: %v", err)
    }

    from := mailerURL.User.Username()
    password, _ := mailerURL.User.Password()
    smtpHost := mailerURL.Hostname()
    smtpPort := mailerURL.Port()

    log.Printf("SMTP Host: %s, SMTP Port: %s, From: %s", smtpHost, smtpPort, from)

    // not working TODO
    msg := "From: " + "forem@test.scrape" + "\n" + // make enything 
        "To: " + to + "\n" + // this should be the user logged
        "Subject: " + subject + "\n\n" +
        "Content-Type: text/html; charset=UTF-8\n\n" +
        body

    auth := smtp.PlainAuth("", from, password, smtpHost)

    tlsConfig := &tls.Config{
        InsecureSkipVerify: true, // 4 test
        ServerName:         smtpHost,
    }

    conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
    if err != nil {
        return fmt.Errorf("failed to dial SMTP server: %v", err)
    }
    defer conn.Close()

    client, err := smtp.NewClient(conn, smtpHost)
    if err != nil {
        return fmt.Errorf("failed to create SMTP client: %v", err)
    }

    if err = client.Auth(auth); err != nil {
        return fmt.Errorf("failed to authenticate: %v", err)
    }

    if err = client.Mail(from); err != nil {
        return fmt.Errorf("failed to set sender: %v", err)
    }

    if err = client.Rcpt(to); err != nil {
        return fmt.Errorf("failed to set recipient: %v", err)
    }

    w, err := client.Data()
    if err != nil {
        return fmt.Errorf("failed to get data writer: %v", err)
    }
    defer w.Close()

    _, err = w.Write([]byte(msg))
    if err != nil {
        return fmt.Errorf("failed to write message: %v", err)
    }

    log.Println("Email sent successfully.")
    return nil
}