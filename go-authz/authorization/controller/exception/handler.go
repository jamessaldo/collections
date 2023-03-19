package exception

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestExceptionHandler(ctx *gin.Context, err BadRequestException) {
	ctx.JSON(http.StatusBadRequest, gin.H{"code": err.Code(), "status": "fail", "detail": err.Error()})
}

func ForbiddenExceptionHandler(ctx *gin.Context, err ForbiddenException) {
	ctx.JSON(http.StatusForbidden, gin.H{"code": err.Code(), "status": "fail", "detail": err.Error()})
}

func NotFoundExceptionHandler(ctx *gin.Context, err NotFoundException) {
	ctx.JSON(http.StatusNotFound, gin.H{"code": err.Code(), "status": "fail", "detail": err.Error()})
}

func UnauthorizedExceptionHandler(ctx *gin.Context, err UnauthorizedException) {
	ctx.JSON(http.StatusUnauthorized, gin.H{"code": err.Code(), "status": "fail", "detail": err.Error()})
}

func BadGatewayExceptionHandler(ctx *gin.Context, err BadGatewayException) {
	ctx.JSON(http.StatusBadGateway, gin.H{"code": err.Code(), "status": "fail", "detail": err.Error()})
}

func ConflictExceptionHandler(ctx *gin.Context, err ConflictException) {
	ctx.JSON(http.StatusConflict, gin.H{"code": err.Code(), "status": "fail", "detail": err.Error()})
}

func DefaultExceptionHandler(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "status": "fail", "detail": err.Error()})
}
