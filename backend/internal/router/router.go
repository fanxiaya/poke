package router

import (
	"poke/backend/internal/handler"
	"poke/backend/internal/validatorx"

	"github.com/gin-gonic/gin"
)

// New 构建并返回 Gin 路由引擎。
func New() *gin.Engine {
	if err := validatorx.RegisterCustomValidators(); err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	healthHandler := handler.NewHealth()
	r.GET("/health", healthHandler.Check)

	api := r.Group("/api")
	{
		_ = api
	}

	return r
}
