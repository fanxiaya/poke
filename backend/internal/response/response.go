package response

import (
	"net/http"

	"poke/backend/internal/errno"

	"github.com/gin-gonic/gin"
)

// Body 定义统一响应结构。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Body{
		Code:    errno.CodeSuccess,
		Message: errno.Message(errno.CodeSuccess),
		Data:    data,
	})
}

// Fail 返回失败响应，可覆盖默认 message。
func Fail(c *gin.Context, code int, message string) {
	if message == "" {
		message = errno.Message(code)
	}
	c.JSON(http.StatusOK, Body{
		Code:    code,
		Message: message,
	})
}
