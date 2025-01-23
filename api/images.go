package api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path"
	"strings"

	db "github.com/Ecc-asplay/backend/db/sqlc"
	"github.com/Ecc-asplay/backend/token"
	"github.com/Ecc-asplay/backend/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) saveImagesForPost(ctx *gin.Context, postID uuid.UUID, files []*multipart.FileHeader) error {
	s3Client, err := util.GetS3Client(ctx, s.config)
	if err != nil {
		return fmt.Errorf("S3クライアントの取得に失敗しました: %w", err)
	}

	bucket := s.config.BucketName
	if bucket == "" {
		return fmt.Errorf("S3バケット名が指定されていません")
	}

	// 最大 5 枚分の画像URL を格納するためのスライス
	var imageURLs [5]string

	for i, fileHeader := range files {
		if i >= 5 {
			break
		}

		file, err := fileHeader.Open()
		if err != nil {
			return fmt.Errorf("ファイルを開けませんでした: %w", err)
		}

		buffer := bytes.NewBuffer(nil)
		if _, err := buffer.ReadFrom(file); err != nil {
			file.Close() // エラー時は即閉じる
			return fmt.Errorf("ファイルデータの読み取りに失敗しました: %w", err)
		}

		file.Close()

		key := "images/" + "image_" + postID.String() + "_" + fileHeader.Filename

		_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String(bucket),
			Key:         aws.String(key),
			Body:        bytes.NewReader(buffer.Bytes()),
			ContentType: aws.String(getContentType(fileHeader.Filename)),
		})
		if err != nil {
			return fmt.Errorf("S3へのファイルのアップロードに失敗しました: %w", err)
		}

		imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)
		imageURLs[i] = imageURL
	}

	_, err = s.store.CreateImage(ctx, db.CreateImageParams{
		PostID: postID,
		Page:   1, // 必要に応じて調整
		Img1:   []byte(imageURLs[0]),
		Img2:   []byte(imageURLs[1]),
		Img3:   []byte(imageURLs[2]),
		Img4:   []byte(imageURLs[3]),
		Img5:   []byte(imageURLs[4]),
	})
	if err != nil {
		return fmt.Errorf("failed to save image URLs to DB: %w", err)
	}

	return nil
}

func (s *Server) UploadImagesHandler(ctx *gin.Context) {

	// 認証トークンの確認（一応）
	if _, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload); !ok {
		ctx.JSON(401, gin.H{"error": "認証に失敗しました"})
		return
	}

	postIDStr := ctx.PostForm("postID")
	if postIDStr == "" {
		ctx.JSON(400, gin.H{"error": "投稿IDが必要です"})
		return
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "無効な投稿IDです"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, gin.H{"error": "不正なフォームデータです"})
		return
	}
	files := form.File["images"]
	if len(files) == 0 {
		ctx.JSON(400, gin.H{"error": "画像ファイルが添付されていません"})
		return
	}

	if err := s.saveImagesForPost(ctx, postID, files); err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("画像アップロードに失敗しました: %v", err)})
		return
	}

	ctx.JSON(200, gin.H{"message": "画像のアップロードに成功しました"})
}

func getContentType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
