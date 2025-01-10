package util

import (
	"fmt"
	"net/smtp"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func SendMail(config Config, to []string, subject, body string) error {
	from := fmt.Sprintf("%s <%s>", config.SmtpFromName, config.SmtpFromAddress)
	msg := "From: " + from + "\n" +
		"Subject: " + subject + "\n" +
		"Content-Type: text/html; charset=utf-8\n\n" +
		body

	// SMTP認証
	auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPassword, config.SmtpHost)

	// メール送信
	err := smtp.SendMail(config.SmtpHost+":"+strconv.Itoa(config.SmtpPort), auth, config.SmtpFromAddress, to, []byte(msg))
	if err != nil {
		return fmt.Errorf("メールの送信に失敗しました: %w", err)
	}
	return nil
}

func SaveVerificationCode(rdb *redis.Client, email string, code string, expiration time.Duration) error {
	err := rdb.Set(email, code, expiration).Err()
	if err != nil {
		return fmt.Errorf("認証コードの保存に失敗しました: %w", err)
	}
	return nil
}
