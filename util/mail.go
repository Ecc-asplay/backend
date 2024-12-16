package util

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/go-redis/redis"
)

type MailConfig struct {
	SMTPHost    string
	SMTPPort    string
	Username    string
	Password    string
	FromAddress string
	FromName    string
}

func SendMail(config MailConfig, to []string, subject, body string) error {
	from := fmt.Sprintf("%s <%s>", config.FromName, config.FromAddress)

	msg := "From: " + from + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		body

	// SMTP認証
	auth := smtp.PlainAuth("", config.Username, config.Password, config.SMTPHost)

	// メール送信
	err := smtp.SendMail(config.SMTPHost+":"+config.SMTPPort, auth, config.FromAddress, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func GenerateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}

// 認証コードを保存
func SaveVerificationCode(rdb *redis.Client, email string, code string, expiration time.Duration) error {
	err := rdb.Set(email, code, expiration).Err()
	if err != nil {
		return fmt.Errorf("認証コードの保存に失敗しました: %w", err)
	}
	return nil
}
