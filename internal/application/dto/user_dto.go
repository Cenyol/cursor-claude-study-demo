package dto

// RegisterRequest 注册请求 DTO
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginRequest 登录请求 DTO
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResult 注册成功返回
type RegisterResult struct {
	User map[string]interface{} `json:"user"`
}

// LoginResult 登录成功返回
type LoginResult struct {
	Token string `json:"token"`
}
