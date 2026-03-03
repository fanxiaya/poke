package handler

import (
	"poke/backend/internal/response"

	"github.com/gin-gonic/gin"
)

// Health 健康检查接口处理器。
type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

// Check 返回服务健康状态。
func (h *Health) Check(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
	})
}
