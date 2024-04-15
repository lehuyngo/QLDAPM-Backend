package http_parser

import (
	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpCode int, data any) {
	c.JSON(httpCode, data)
}
