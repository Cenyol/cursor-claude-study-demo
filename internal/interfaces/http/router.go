package http

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter 注册路由
func SetupRouter(h *UserHandler) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", h.Register)
		v1.POST("/login", h.Login)
		v1.POST("/logout", h.Logout)
	}
	return r
}
