package error

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrPermissionDenied = errors.New("permission denied")
	ErrConflict         = errors.New("conflict")
)

func HandleDBError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidInput):
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid input", 400))
	case errors.Is(err, ErrPermissionDenied):
		ctx.JSON(http.StatusForbidden, errorResponse("permission denied", 403))
	case errors.Is(err, sql.ErrNoRows):
		ctx.JSON(http.StatusNotFound, errorResponse("record not found", 404))
	case errors.Is(err, ErrConflict):
		ctx.JSON(http.StatusConflict, errorResponse("resource conflict", 409))
	default:
		ctx.JSON(http.StatusInternalServerError, errorResponse("internal server error", 500))
	}
}

func errorResponse(message string, code int) gin.H {
	return gin.H{
		"Error": message,
		"Code":  code,
	}
}
