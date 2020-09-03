package proto

// GuestLoginReq 游客登录请求
// :port/v1/login/guest
type GuestLoginReq struct {
}

// GuestLoginRes 游客登录返回
type GuestLoginRes struct {
	Token string
}

// LoginOutReq 登出
type LoginOutReq struct {
	Token string
}
