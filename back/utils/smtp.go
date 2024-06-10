package utils

import (
    "crypto/tls"
    "encoding/base64"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "mime/quotedprintable"
    "net/smtp"
    "net/url"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

func LoadEnv() error {
    return godotenv.Load()
}

func SendEmail(to string, subject string, body string, attachmentPath string) error {
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

    // Leer el archivo CSV y codificarlo en base64
    csvData, err := ioutil.ReadFile(attachmentPath)
    if err != nil {
        return fmt.Errorf("failed to read attachment: %v", err)
    }
    encodedAttachment := base64.StdEncoding.EncodeToString(csvData)

    boundary := "my-boundary-1234567890"

    // Construir el mensaje de email con el adjunto
    var msg strings.Builder
    msg.WriteString("From: " + "scrapping@forem" + "\n")
    msg.WriteString("To: " + to + "\n")
    msg.WriteString("Subject: " + subject + "\n")
    msg.WriteString("MIME-Version: 1.0\n")
    msg.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\n\n")

    // Parte del cuerpo del email
    msg.WriteString("--" + boundary + "\n")
    msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\n")
    msg.WriteString("Content-Transfer-Encoding: quoted-printable\n\n")
    writeQuotedPrintable(&msg, body)
    msg.WriteString("\n\n")

    // Parte del adjunto
    msg.WriteString("--" + boundary + "\n")
    msg.WriteString("Content-Type: text/csv; name=\"" + attachmentPath + "\"\n")
    msg.WriteString("Content-Transfer-Encoding: base64\n")
    msg.WriteString("Content-Disposition: attachment; filename=\"" + attachmentPath + "\"\n\n")
    msg.WriteString(encodedAttachment)
    msg.WriteString("\n\n")
    msg.WriteString("--" + boundary + "--")

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

    _, err = w.Write([]byte(msg.String()))
    if err != nil {
        return fmt.Errorf("failed to write message: %v", err)
    }

    log.Println("Email sent successfully.")
    return nil
}

func writeQuotedPrintable(w io.Writer, body string) {
    qpw := quotedprintable.NewWriter(w)
    defer qpw.Close()
    _, _ = qpw.Write([]byte(body))
}