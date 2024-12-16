package util

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
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

// 環境変数からMailConfigをロード
func LoadMailConfig() MailConfig {
	config := MailConfig{
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    os.Getenv("SMTP_PORT"),
		Username:    os.Getenv("SMTP_USER"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		FromAddress: os.Getenv("SMTP_FROM_ADDRESS"),
		FromName:    os.Getenv("SMTP_FROMNAME"),
	}

	fmt.Printf("Loaded MailConfig: %+v\n", config) // デバッグ出力

	return config
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
