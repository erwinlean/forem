package utils

import (
    "crypto/tls"
    "encoding/base64"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "mime/multipart"
    "mime/quotedprintable"
    "net/smtp"
    "net/textproto"
    "net/url"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

func LoadEnv() error {
    return godotenv.Load()
}

func SendEmail(to string, subject string, body string, attachmentPath string, newAttachmentPath string) error {
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

    csvData, err := ioutil.ReadFile(attachmentPath)
    if err != nil {
        return fmt.Errorf("failed to read attachment: %v", err)
    }

    newCsvData, err := ioutil.ReadFile(newAttachmentPath)
    if err != nil {
        return fmt.Errorf("failed to read new attachment: %v", err)
    }

    // build
    var msg strings.Builder
    msg.WriteString("From: " + "forem@scrapping" + "\n")
    msg.WriteString("To: " + to + "\n")
    msg.WriteString("Subject: " + subject + "\n")
    msg.WriteString("MIME-Version: 1.0\n")

    writer := multipart.NewWriter(&msg)
    boundary := writer.Boundary()
    msg.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n\n", boundary))

    bodyHeader := make(textproto.MIMEHeader)
    bodyHeader.Set("Content-Type", "text/html; charset=\"UTF-8\"")
    bodyHeader.Set("Content-Transfer-Encoding", "quoted-printable")
    bodyPart, err := writer.CreatePart(bodyHeader)
    if err != nil {
        return fmt.Errorf("failed to create body part: %v", err)
    }
    writeQuotedPrintable(bodyPart, body)

    attachmentHeader := make(textproto.MIMEHeader)
    attachmentHeader.Set("Content-Type", "text/csv")
    attachmentHeader.Set("Content-Transfer-Encoding", "base64")
    attachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", attachmentPath))
    attachmentPart, err := writer.CreatePart(attachmentHeader)
    if err != nil {
        return fmt.Errorf("failed to create attachment part: %v", err)
    }
    encodeAndWriteBase64(attachmentPart, csvData)

    newAttachmentHeader := make(textproto.MIMEHeader)
    newAttachmentHeader.Set("Content-Type", "text/csv")
    newAttachmentHeader.Set("Content-Transfer-Encoding", "base64")
    newAttachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", newAttachmentPath))
    newAttachmentPart, err := writer.CreatePart(newAttachmentHeader)
    if err != nil {
        return fmt.Errorf("failed to create new attachment part: %v", err)
    }
    encodeAndWriteBase64(newAttachmentPart, newCsvData)

    err = writer.Close()
    if err != nil {
        return fmt.Errorf("failed to close writer: %v", err)
    }

    auth := smtp.PlainAuth("", from, password, smtpHost)

    tlsConfig := &tls.Config{
        InsecureSkipVerify: true, // test
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

func encodeAndWriteBase64(w io.Writer, data []byte) {
    b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
    base64.StdEncoding.Encode(b, data)
    w.Write(b)
}