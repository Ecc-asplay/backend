package api

import (
	"net/http"
	"time"

	"github.com/Ecc-asplay/backend/util"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func (s *Server) SendVerificationEmail(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な入力データです"})
		return
	}

	// 環境変数からMailConfigをロード
	mailConfig := util.LoadMailConfig()

	// 認証コードを生成
	verificationCode := util.GenerateVerificationCode()

	// 認証コードを保存（有効期限5分）
	err := s.redis.Set(req.Email, verificationCode, 5*time.Minute).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "認証コードの保存に失敗しました", "details": err.Error()})
		return
	}

	// メール内容
	subject := "アカウント認証コードのお知らせ"
	body := `<html>
		<body>
			<h3>「やみよあけ」をご利用いただき、誠にありがとうございます。
			アカウントのセキュリティを確保するため、確認コードをお送りいたします。</h3>
			<p>認証コード:<strong>` + verificationCode + `</strong></p>
			<h3>このコードは、ウェブサイトまたはアプリの指定された欄に、5分以内に入力してアカウント認証を完了してください。</h3>
			<h3>【重要】</h3>
			<ul>
			<li>このメールに心当たりのない場合は、すぐにサポートチームまでお問い合わせください。</li>
			<li>このコードを他人に教えないでください。やみよあけの担当者がお客様にコードを尋ねることは絶対にありません。</li>
			</ul>
			<h3>ご不明な点がございましたら、お気軽にサポートチームまでお問い合わせください。</h3>
			<br>
			<h2>敬具</h2>
			<h2>やみよあけ チーム</h2>
			<p>※このメールは送信専用です。返信はしないでください。</p>
		</body>
    </html>`

	// メール送信
	err = util.SendMail(mailConfig, []string{req.Email}, subject, body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "メール送信に失敗しました", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "認証メールが送信されました"})
}

// 認証コード確認エンドポイント
func (s *Server) VerifyCode(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な入力データです"})
		return
	}

	// 認証コードを取得
	storedCode, err := s.redis.Get(req.Email).Result()
	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "認証コードが見つかりません"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "認証コードの取得に失敗しました", "details": err.Error()})
		return
	}

	// ユーザーが入力したコードと比較
	if storedCode != req.Code {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証コードが無効です"})
		return
	}

	// 認証コードが一致した場合、Redis から削除する（オプション）
	err = s.redis.Del(req.Email).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "認証コードの削除に失敗しました", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "認証コードが確認されました"})
}
