package vo

// 登陆+token // 



type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

type CheckTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type CheckTokenResponse struct {
	UserId    uint   `json:"user_id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	ExpiresIn int64  `json:"expires_in"`
}
