package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-uuid"
)

const TraceIDKey = "_clover_trace"

func Trace(c *gin.Context) {

	traceId, _ := uuid.GenerateUUID()
	c.Set(TraceIDKey, traceId)
	c.Next()
}
