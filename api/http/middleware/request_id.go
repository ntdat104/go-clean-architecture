package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-clean-architecture/pkg/uuid"
)

const (
	// RequestIDHeader is the header key for request ID
	RequestIDHeader = "X-Request-ID"
)

// RequestID is a middleware that injects a request ID into the context of each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID from header or generate a new one
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewGoogleUUID()
		}

		// Set request ID to header
		c.Writer.Header().Set(RequestIDHeader, requestID)
		c.Set(RequestIDHeader, requestID)

		c.Next()
	}
}
