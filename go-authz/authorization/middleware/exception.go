package middleware

import (
	"auth/controller/exception"

	"github.com/gin-gonic/gin"
)

func HandleCustomError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if len(ctx.Errors) > 0 {
				err := ctx.Errors[0].Err

				switch e := err.(type) {
				case exception.BadRequestException:
					exception.BadRequestExceptionHandler(ctx, e)
				case exception.ForbiddenException:
					exception.ForbiddenExceptionHandler(ctx, e)
				case exception.NotFoundException:
					exception.NotFoundExceptionHandler(ctx, e)
				case exception.UnauthorizedException:
					exception.UnauthorizedExceptionHandler(ctx, e)
				case exception.BadGatewayException:
					exception.BadGatewayExceptionHandler(ctx, e)
				case exception.ConflictException:
					exception.ConflictExceptionHandler(ctx, e)
				default:
					exception.DefaultExceptionHandler(ctx, e)
				}

				ctx.Abort()
			}
		}()

		ctx.Next()
	}
}
