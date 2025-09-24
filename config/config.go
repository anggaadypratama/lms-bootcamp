package config

import (
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Port        string `env:"PORT,required"`
	JWTSecret   string `env:"JWT_SECRET,required"`
	Mailer      string `env:"MAIL_MAILER,required"`
	MailHost    string `env:"MAIL_HOST,required"`
	MailPort    string `env:"MAIL_PORT,required"`
	MailUser    string `env:"MAIL_USER,required"`
	MailPass    string `env:"MAIL_PASS,required"`
	MailEncrypt string `env:"MAIL_ENCRYPTION,required"`
	MailFrom    string `env:"MAIL_FROM_ADDRESS,required"`
	MailName    string `env:"MAIL_FROM_NAME,required"`
	AWS_ACCESS_KEY_ID     	string `env:"AWS_ACCESS_KEY_ID,required"`
	AWS_SECRET_ACCESS_KEY  	string `env:"AWS_SECRET_ACCESS_KEY,required"`
	AWS_REGION             	string `env:"AWS_REGION,required"`
	AWS_ENDPOINT           	string `env:"AWS_ENDPOINT,required"`
	AWS_BUCKET            	string `env:"AWS_BUCKET,required"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Mailer: 	os.Getenv("MAIL_MAILER"),
		MailHost: 	os.Getenv("MAIL_HOST"),
		MailPort: 	os.Getenv("MAIL_PORT"),
		MailUser: 	os.Getenv("MAIL_USER"),
		MailPass: 	os.Getenv("MAIL_PASS"),
		MailEncrypt: os.Getenv("MAIL_ENCRYPTION"),
		MailFrom: 	os.Getenv("MAIL_FROM_ADDRESS"),
		MailName: 	os.Getenv("MAIL_FROM_NAME"),
		AWS_ACCESS_KEY_ID:      os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_REGION:              os.Getenv("AWS_REGION"),
		AWS_ENDPOINT:            os.Getenv("AWS_ENDPOINT"),
		AWS_BUCKET:              os.Getenv("AWS_BUCKET"),
	}

	return cfg, nil
}
