package dto

// WxLoginRequest 微信小程序登录请求。
type WxLoginRequest struct {
	Code     string `json:"code" binding:"required,min=1,max=128"`
	Nickname string `json:"nickname" binding:"max=50"`
	AvatarURL string `json:"avatarUrl" binding:"omitempty,url,max=512"`
}

