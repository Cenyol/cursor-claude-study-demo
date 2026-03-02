package http

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"user-system/internal/application/dto"
	"user-system/internal/application/usecase"
	"user-system/internal/domain/repository"
)

// UserHandler 接口层：仅处理 HTTP 与参数校验，编排调用 Application 用例
type UserHandler struct {
	UserRepo    repository.UserRepository
	SessionRepo repository.SessionRepository
}

// Register 注册
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONBadRequest(c, "invalid params")
		return
	}
	if req.Username == "" || req.Password == "" {
		JSONBadRequest(c, "username and password required")
		return
	}
	result, err := usecase.Register(c.Request.Context(), &req, h.UserRepo)
	if err != nil {
		if err == usecase.ErrUsernameExists {
			JSONConflict(c, "username already exists")
			return
		}
		if err == usecase.ErrPasswordTooShort {
			JSONBadRequest(c, "password at least 8 characters")
			return
		}
		JSONServerError(c, "register failed")
		return
	}
	JSONSuccess(c, result)
}

// Login 登录
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		JSONBadRequest(c, "invalid params")
		return
	}
	if req.Username == "" || req.Password == "" {
		JSONBadRequest(c, "username and password required")
		return
	}
	result, err := usecase.Login(c.Request.Context(), &req, h.UserRepo, h.SessionRepo)
	if err != nil {
		if err == usecase.ErrUserNotFound || err == usecase.ErrInvalidPassword {
			JSONUnauthorized(c, "invalid username or password")
			return
		}
		JSONServerError(c, "login failed")
		return
	}
	JSONSuccess(c, result)
}

// Logout 退出（Authorization: Bearer <token>）
func (h *UserHandler) Logout(c *gin.Context) {
	tok := extractBearerToken(c)
	if tok == "" {
		JSONUnauthorized(c, "missing token")
		return
	}
	if err := usecase.Logout(context.Background(), tok, h.UserRepo, h.SessionRepo); err != nil {
		JSONServerError(c, "logout failed")
		return
	}
	JSONSuccess(c, nil)
}

func extractBearerToken(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return ""
	}
	return strings.TrimSpace(auth[len(prefix):])
}
