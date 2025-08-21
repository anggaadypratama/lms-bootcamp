package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"strconv"

	"lms-bootcamp/config"

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
    cfg    *config.Config
}

func NewEmailService() *EmailService {
    conf, err := config.NewConfig()
    if err != nil {
        log.Fatal("Failed to load config:", err)
    }

	port, err := strconv.Atoi(conf.MailPort)
	if err != nil {
		log.Fatal("Invalid MAIL_PORT:", err)
	}

	dialer := gomail.NewDialer(conf.MailHost, port, conf.MailUser, conf.MailPass)
    dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailService{
		Dialer: dialer,
		cfg:    conf,
	}
}

func (s *EmailService) SendEmailWithChannel(req EmailRequest) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                req.ResultChan <- fmt.Errorf("panic: %v", r)
            }
            close(req.ResultChan)
        }()

        msg := gomail.NewMessage()
        msg.SetHeader("From", s.cfg.MailFrom)
        msg.SetHeader("To", req.To)
        msg.SetHeader("Subject", req.Subject)
        msg.SetBody("text/plain", req.Body)
        
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