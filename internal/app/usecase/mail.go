package usecase

import (
	"crypto/tls"
	"file_storage_service/infrastructure/config"
	"file_storage_service/internal/domain/model"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type MailUsecase struct {
	Config config.MailConfig
}

func NewMailUsecase(cfg config.MailConfig) *MailUsecase {
	return &MailUsecase{Config: cfg}
}

func (u *MailUsecase) SendMail(req model.MailRequest) error {
	host := u.Config.Host
	port := u.Config.Port
	from := u.Config.From
	user := u.Config.Username
	pass := u.Config.Password
	insecureSkipVerify := u.Config.InsecureSkipVerify
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("Mail server address:", addr)

	auth := smtp.PlainAuth("", user, pass, host)

	// MIME (รองรับ HTML)
	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", req.To))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", req.Subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(req.Body)

	// STARTTLS สำหรับพอร์ต 587
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("dial failed: %w", err)
	}
	defer client.Quit()

	// สั่งเริ่ม TLS
	tlsconfig := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify, // สำหรับ Mailtrap เท่านั้น
		ServerName:         host,
	}
	if err = client.StartTLS(tlsconfig); err != nil {
		return fmt.Errorf("starttls failed: %w", err)
	}

	// Auth
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	// ระบุผู้ส่งและผู้รับ
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("mail from failed: %w", err)
	}
	if err = client.Rcpt(req.To); err != nil {
		return fmt.Errorf("rcpt failed: %w", err)
	}

	// เขียน message และส่ง
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data failed: %w", err)
	}

	_, err = w.Write([]byte(message.String()))
	if err != nil {
		return fmt.Errorf("write message failed: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("close failed: %w", err)
	}

	log.Printf("✅ ส่งอีเมลไปยัง %s เรียบร้อย", req.To)
	return nil
}
