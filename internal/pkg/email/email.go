package email

import (
	"crypto/tls"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailRequest struct {
    To      string
    Subject string
    Body    string
    ResultChan chan error
}

type EmailService struct {
	Dialer *gomail.Dialer
}

func NewEmailService() *EmailService {
	dialer := gomail.NewDialer("localhost", 1025, "", "")
    dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailService{
		Dialer: dialer,
	}
}

func (s *EmailService) SendEmailWithChannel(req EmailRequest) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                fmt.Println("Panic saat kirim email:", r)
                req.ResultChan <- fmt.Errorf("panic: %v", r)
            }
            close(req.ResultChan)
        }()

        msg := gomail.NewMessage()
        msg.SetHeader("From", "shilla@beloved.com")
        msg.SetHeader("To", req.To)
        msg.SetHeader("Subject", req.Subject)
        msg.SetBody("text/plain", req.Body)

        fmt.Println("Connecting to SMTP server...")
        err := s.Dialer.DialAndSend(msg)
        if err != nil {
            log.Println("Gagal kirim email:", err)
        } else {
            log.Println("Email berhasil dikirim ke:", req.To)
        }
        req.ResultChan <- err
    }()
}

func (s *EmailService) SendEmailAsyncWithResult(to, subject, body string) <-chan error {
    resultChan := make(chan error, 1)
    
    req := EmailRequest{
        To:         to,
        Subject:    subject,
        Body:       body,
        ResultChan: resultChan,
    }
    
    s.SendEmailWithChannel(req)
    return resultChan
}