package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	db "github.com/Ecc-asplay/backend/db/sqlc"
)

type CreatePostRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	PostID   uuid.UUID `json:"post_id"`
	ShowID   string    `json:"show_id"`
	Title    string    `json:"title"`
	Feel     string    `json:"feel"`
	Content  []byte    `json:"content"`
	Reaction int32     `json:"reaction"`
	Status   string    `json:"status"`
}

func (s *Server) CreatePost(ctx *gin.Context) {
	var req CreatePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	postData := db.CreatePostParams{
		UserID:   req.UserID,
		PostID:   req.PostID,
		ShowID:   req.ShowID,
		Title:    req.Title,
		Feel:     req.Feel,
		Content:  req.Content,
		Reaction: req.Reaction,
		Status:   req.Status,
	}

	post, err := s.store.CreatePost(ctx, postData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// Get
func (s *Server) GetAllPost(ctx *gin.Context) {
	post, err := s.store.GetPostsList(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// var imageData2 []map[string]interface{}
	// err1 := json.Unmarshal(post[0].Images, &imageData2)
	// if err1 != nil {
	// 	log.Fatalf("Error decoding JSON: %v", err1)
	// }

	// var img [][][]byte // [] <= page , [] <= img number , []byte

	// for _, pageData := range imageData2 {
	// 	page := pageData["page"].(string)
	// 	image1 := pageData["image1"].([]byte)
	// 	image2 := pageData["image2"].([]byte)
	// 	image3 := pageData["image3"].([]byte)
	// 	image4 := pageData["image4"].([]byte)
	// 	image5 := pageData["image5"].([]byte)

	// 	pageIndex := -1
	// 	switch page {
	// 	case "1":
	// 		pageIndex = 0
	// 	case "2":
	// 		pageIndex = 1
	// 	case "3":
	// 		pageIndex = 2
	// 	}
	// 	if pageIndex >= 0 {
	// 		if len(img) <= pageIndex {
	// 			img = append(img, make([][]byte, 5))
	// 		}
	// 		img[pageIndex][0] = image1
	// 		img[pageIndex][1] = image2
	// 		img[pageIndex][2] = image3
	// 		img[pageIndex][3] = image4
	// 		img[pageIndex][4] = image5
	// 	}
	// }
	// log.Println(post[1].CreatedAt)

	ctx.JSON(http.StatusOK, post)
}

func (s *Server) GetPostOfKeywords(ctx *gin.Context) {
	var req string
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	findPost, err := s.store.GetPostOfKeywords(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, findPost)
}

// Delete

type DeletePostRequest struct {
	UserID uuid.UUID `json:"user_id"`
	PostID uuid.UUID `json:"post_id"`
}

func (s *Server) DeletePost(ctx *gin.Context) {
	var req DeletePostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := s.store.DeletePost(ctx, db.DeletePostParams{
		UserID: req.UserID,
		PostID: req.PostID,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.Status(http.StatusOK)
}
