package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"github.com/Ecc-asplay/backend/util"
)

type SendVerificationEmailReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (s *Server) SendVerificationEmail(ctx *gin.Context) {
	var req SendVerificationEmailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "無効な入力データです")
		return
	}

	mailConfig := util.LoadMailConfig(s.config)
	// 認証コードを生成
	verificationCode := util.GenerateVerificationCode()

	// 認証コードを保存（有効期限5分）
	err := s.redis.Set(req.Email, verificationCode, 5*time.Minute).Err()
	if err != nil {
		handleDBError(ctx, err, "認証コードの保存に失敗しました")
		return
	}

	// メール内容
	subject := "「やみよあけ」認証コードのご案内"
	body := `<html>
		<body>
			<h3>「やみよあけ」をご利用いただき、ありがとうございます。</h3>
			<p>アカウントのセキュリティを確保するため、以下の認証コードをお送りします。</p>
			<p><strong>認証コード：` + verificationCode + `</strong></p>
			<p>このコードは、ウェブサイトまたはアプリの指定された欄に、<strong>5分以内に入力</strong>して認証を完了してください。</p>
			
			<h3>【ご注意ください】</h3>
			<ul>
				<li>このメールに心当たりがない場合は、すぐにサポートチームまでご連絡ください。</li>
				<li>認証コードを他人に教えないでください。「やみよあけ」の担当者が認証コードを尋ねることは一切ありません。</li>
			</ul>
	
			<p>ご不明な点がございましたら、以下のリンクからサポートチームまでお問い合わせください。</p>
			<p><a href="mailto:support@yamiyoake.com">support@yamiyoake.com</a></p>
			
			<p style="font-size:small;">※このメールは送信専用です。返信しないようお願いいたします。</p>
		</body>
	</html>`

	// メール送信
	err = util.SendMail(mailConfig, []string{req.Email}, subject, body)
	if err != nil {
		handleDBError(ctx, err, "メール送信に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "認証メールが送信されました", "email": req.Email})
}

type VerifyCodeReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

// 認証コード確認
func (s *Server) VerifyCode(ctx *gin.Context) {
	var req VerifyCodeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleDBError(ctx, err, "無効な入力データです")
		return
	}

	// 認証コードを取得
	storedCode, err := s.redis.Get(req.Email).Result()
	if err != nil {
		if err == redis.Nil {
			handleDBError(ctx, err, "認証コードが見つかりません")
			return
		}
		handleDBError(ctx, err, "認証コードの取得に失敗しました")
		return
	}

	// ユーザーが入力したコードと比較
	if storedCode != req.Code {
		handleDBError(ctx, err, "認証コードが無効です")
		return
	}

	// 認証コードが一致した場合、Redis から削除する
	err = s.redis.Del(req.Email).Err()
	if err != nil {
		handleDBError(ctx, err, "認証コードの削除に失敗しました")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "認証コードが確認されました"})
}
