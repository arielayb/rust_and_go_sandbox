package app

const (
	USER_INFO = "USER_INFO"
)

type UserWebInfo struct {
	UserID  string `json:"user_id"`
	Method  string `json:"method"`
	Message string `json:"msg"`
	Global  bool   `json:"global"`
}
