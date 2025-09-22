package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/api/http/app_context"
	"github.com/ntdat104/go-clean-architecture/pkg/logger"
	"go.uber.org/zap"
)

func AppContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := &app_context.AppContext{
			Ctx:    c.Request.Context(),
			Logger: logger.Logger.With(zap.String("request_id", c.GetString("X-Request-ID"))),
		}
		defer appCtx.Cleanup()
		c.Set("app_context", appCtx)
		c.Next()
	}
}
