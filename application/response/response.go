package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/pkg/uuid"
)

type Meta struct {
	MessageID string `json:"message_id"`
	Timestamp int64  `json:"timestamp"`
	Datetime  string `json:"datetime"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Token     string `json:"token,omitempty"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data,omitempty"`
}

func getMessageID(ctx *gin.Context) string {
	messageID := ctx.GetHeader("X-Message-ID")
	if messageID == "" {
		messageID = uuid.NewShortUUID()
	}
	return messageID
}

func buildResponse(ctx *gin.Context, code int, obj any) {
	now := time.Now()
	response := Response{
		Meta: Meta{
			MessageID: getMessageID(ctx),
			Timestamp: now.UnixMilli(),
			Datetime:  now.Format("2006-01-02 15:04:05"),
			Code:      code,
			Message:   http.StatusText(code),
		},
		Data: obj,
	}
	ctx.Header("X-Message-ID", getMessageID(ctx))
	ctx.JSON(code, response)
}

func JSON(ctx *gin.Context, code int, obj any) {
	buildResponse(ctx, code, obj)
}

func Success(ctx *gin.Context, obj any) {
	buildResponse(ctx, http.StatusOK, obj)
}

func Failed(ctx *gin.Context, obj any) {
	buildResponse(ctx, http.StatusBadRequest, obj)
}
